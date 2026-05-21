package license

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRevokeLicenseTreatsStatusNotFoundAsSoftFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	svc := NewService("", "", "", server.URL, "", "", "")
	if err := svc.RevokeLicense(context.Background(), "lic-1"); err != nil {
		t.Fatalf("expected 404 to be tolerated, got %v", err)
	}
}

func TestRevokeLicenseStillFailsForServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer server.Close()

	svc := NewService("", "", "", server.URL, "", "", "")
	if err := svc.RevokeLicense(context.Background(), "lic-1"); err == nil {
		t.Fatal("expected server error to fail")
	}
}
