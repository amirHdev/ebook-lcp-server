package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginReturnsAdminTokenWithTwoFactor(t *testing.T) {
	handler := NewAuthHandler("secret", "toghyani", "admin-pass", "publisher01", "publisher-pass", "65b5ec", "tenant-a")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"username":"toghyani","password":"admin-pass","twoFactor":"65b5ec"}`))
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status %d: %s", rec.Code, rec.Body.String())
	}
	var resp LoginResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Role != "admin" || resp.Subject != "toghyani" || resp.Token == "" {
		t.Fatalf("unexpected login response: %#v", resp)
	}
}

func TestLoginReturnsPublisherTokenWithoutTwoFactor(t *testing.T) {
	handler := NewAuthHandler("secret", "toghyani", "admin-pass", "publisher01", "publisher-pass", "65b5ec", "tenant-a")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"username":"publisher01","password":"publisher-pass"}`))
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status %d: %s", rec.Code, rec.Body.String())
	}
	var resp LoginResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Role != "publisher" || resp.Subject != "publisher01" || resp.Token == "" {
		t.Fatalf("unexpected login response: %#v", resp)
	}
}
