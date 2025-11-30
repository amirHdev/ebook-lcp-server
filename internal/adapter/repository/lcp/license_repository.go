package lcp

import (
	"context"
	"sync"

	"github.com/Mehrbod2002/lcp/internal/domain/lcp"
)

type LicenseRepository interface {
	Save(ctx context.Context, license *lcp.License) error
	FindByPublication(ctx context.Context, publicationID *string) ([]*lcp.License, error)
}

type licenseRepository struct {
	mu       sync.RWMutex
	licenses []*lcp.License
}

func NewLicenseRepository() LicenseRepository {
	return &licenseRepository{}
}

func (r *licenseRepository) Save(ctx context.Context, license *lcp.License) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.licenses = append(r.licenses, license)
	return nil
}

func (r *licenseRepository) FindByPublication(ctx context.Context, publicationID *string) ([]*lcp.License, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*lcp.License
	for _, lic := range r.licenses {
		if publicationID == nil || lic.PublicationID == *publicationID {
			result = append(result, lic)
		}
	}
	return result, nil
}
