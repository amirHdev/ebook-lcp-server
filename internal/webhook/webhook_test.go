package webhook

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPPublisherSendsSignedEvent(t *testing.T) {
	var signature string
	var body []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signature = r.Header.Get("X-LCP-Signature")
		body, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	publisher := NewHTTPPublisher([]string{server.URL}, "secret")
	err := publisher.Publish(context.Background(), Event{
		Type:      EventLicenseCreated,
		CreatedAt: time.Unix(1, 0).UTC(),
		Data:      map[string]string{"id": "lic1"},
	})
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}
	if signature == "" {
		t.Fatal("expected signature header")
	}
	if len(body) == 0 {
		t.Fatal("expected request body")
	}
}
