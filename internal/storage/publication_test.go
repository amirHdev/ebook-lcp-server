package storage

import (
	"context"
	"testing"
	"time"
)

func TestParseS3URI(t *testing.T) {
	bucket, key, err := parseS3URI("s3://books/publications/book.epub")
	if err != nil {
		t.Fatalf("parseS3URI failed: %v", err)
	}
	if bucket != "books" || key != "publications/book.epub" {
		t.Fatalf("unexpected result: bucket=%s key=%s", bucket, key)
	}
}

func TestParseS3URIRejectsInvalidValue(t *testing.T) {
	if _, _, err := parseS3URI("https://example.test/book.epub"); err == nil {
		t.Fatal("expected invalid scheme to fail")
	}
}

func TestFilesystemStorageDoesNotSignURLs(t *testing.T) {
	url, ok, err := NewFilesystemPublicationStorage().SignedURL(context.Background(), "/tmp/book.epub", time.Minute)
	if err != nil {
		t.Fatalf("SignedURL failed: %v", err)
	}
	if ok || url != "" {
		t.Fatalf("expected no signed URL, got ok=%v url=%q", ok, url)
	}
}
