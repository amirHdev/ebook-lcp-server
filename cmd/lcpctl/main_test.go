package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHealthCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/healthz":
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		case "/readyz":
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	var out strings.Builder
	if err := run([]string{"health", "--base-url", server.URL}, &out, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"status": "ok"`) {
		t.Fatalf("unexpected output: %s", out.String())
	}
	if !strings.Contains(out.String(), `"status": "ready"`) {
		t.Fatalf("unexpected output: %s", out.String())
	}
}

func TestDemoCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/auth/login":
			_ = json.NewEncoder(w).Encode(map[string]string{"token": "token", "role": "admin"})
		case "/graphql":
			raw, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var payload struct {
				Query     string                 `json:"query"`
				Variables map[string]interface{} `json:"variables"`
			}
			if err := json.Unmarshal(raw, &payload); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			switch {
			case strings.Contains(payload.Query, "uploadPublication"):
				if got := payload.Variables["title"]; got != "Pride and Prejudice" {
					t.Fatalf("unexpected title: %#v", got)
				}
				if _, ok := payload.Variables["file"].(string); !ok {
					t.Fatalf("expected base64 file payload, got %#v", payload.Variables["file"])
				}
				_ = json.NewEncoder(w).Encode(map[string]any{
					"data": map[string]any{
						"uploadPublication": map[string]any{
							"id":    "pub-1",
							"title": "book",
						},
					},
				})
			case strings.Contains(payload.Query, "createLicense"):
				if got := payload.Variables["publicationID"]; got != "pub-1" {
					t.Fatalf("unexpected publicationID: %#v", got)
				}
				if got := payload.Variables["userID"]; got != "reader-01" {
					t.Fatalf("unexpected userID: %#v", got)
				}
				_ = json.NewEncoder(w).Encode(map[string]any{
					"data": map[string]any{
						"createLicense": map[string]any{
							"id":            "lic-1",
							"publicationID": "pub-1",
							"userID":        "reader-01",
						},
					},
				})
			default:
				http.Error(w, "unexpected graphql payload", http.StatusBadRequest)
			}
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	dir := t.TempDir()
	book := filepath.Join(dir, "book.epub")
	if err := os.WriteFile(book, []byte("book"), 0o600); err != nil {
		t.Fatal(err)
	}

	var out strings.Builder
	if err := run([]string{
		"demo",
		"--base-url", server.URL,
		"--file", book,
		"--username", "admin",
		"--password", "admin",
		"--two-factor", "123456",
	}, &out, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"publication_id": "pub-1"`) {
		t.Fatalf("unexpected output: %s", out.String())
	}
	if !strings.Contains(out.String(), `"license_id": "lic-1"`) {
		t.Fatalf("unexpected output: %s", out.String())
	}
}

func TestRevokeLicenseCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/auth/login":
			_ = json.NewEncoder(w).Encode(map[string]string{"token": "token", "role": "admin"})
		case "/graphql":
			raw, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var payload struct {
				Query     string                 `json:"query"`
				Variables map[string]interface{} `json:"variables"`
			}
			if err := json.Unmarshal(raw, &payload); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if !strings.Contains(payload.Query, "revokeLicense") {
				t.Fatalf("unexpected query: %s", payload.Query)
			}
			if got := payload.Variables["id"]; got != "lic-1" {
				t.Fatalf("unexpected license id: %#v", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"revokeLicense": true},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	var out strings.Builder
	if err := run([]string{
		"license",
		"revoke",
		"--base-url", server.URL,
		"--username", "admin",
		"--password", "admin",
		"--two-factor", "123456",
		"--id", "lic-1",
	}, &out, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"status": "revoked"`) {
		t.Fatalf("unexpected output: %s", out.String())
	}
}
