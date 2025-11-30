package lcp

import "context"

// PublicationRepository describes the persistence operations for publications.
type PublicationRepository interface {
	Save(ctx context.Context, pub *Publication) error
	FindAll(ctx context.Context) ([]*Publication, error)
	FindByID(ctx context.Context, id string) (*Publication, error)
}

// LicenseRepository describes the persistence operations for licenses.
type LicenseRepository interface {
	Save(ctx context.Context, license *License) error
	FindByPublication(ctx context.Context, publicationID *string) ([]*License, error)
}
