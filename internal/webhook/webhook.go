package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	EventPublicationUploaded = "publication.uploaded"
	EventLicenseCreated      = "license.created"
	EventLicenseRevoked      = "license.revoked"
)

type Event struct {
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
	Data      any       `json:"data"`
}

type Publisher interface {
	Publish(ctx context.Context, event Event) error
}

type NopPublisher struct{}

func (NopPublisher) Publish(context.Context, Event) error {
	return nil
}

type HTTPPublisher struct {
	urls   []string
	secret string
	client *http.Client
}

func NewHTTPPublisher(urls []string, secret string) Publisher {
	if len(urls) == 0 {
		return NopPublisher{}
	}
	return &HTTPPublisher{
		urls:   append([]string(nil), urls...),
		secret: secret,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (p *HTTPPublisher) Publish(ctx context.Context, event Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	for _, target := range p.urls {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		if p.secret != "" {
			req.Header.Set("X-LCP-Signature", sign(body, p.secret))
		}
		resp, err := p.client.Do(req)
		if err != nil {
			return err
		}
		if err := resp.Body.Close(); err != nil {
			return err
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("webhook returned %s", resp.Status)
		}
	}
	return nil
}

func sign(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}
