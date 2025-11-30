package lcp

import (
	"context"
	"sync"

	"github.com/Mehrbod2002/lcp/internal/domain/lcp"
)

type PublicationRepository interface {
	Save(ctx context.Context, pub *lcp.Publication) error
	FindAll(ctx context.Context) ([]*lcp.Publication, error)
	FindByID(ctx context.Context, id string) (*lcp.Publication, error)
}

type publicationRepository struct {
	mu           sync.RWMutex
	publications []*lcp.Publication
}

func NewPublicationRepository() PublicationRepository {
	return &publicationRepository{}
}

func (r *publicationRepository) Save(ctx context.Context, pub *lcp.Publication) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.publications = append(r.publications, pub)
	return nil
}

func (r *publicationRepository) FindAll(ctx context.Context) ([]*lcp.Publication, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	pubs := make([]*lcp.Publication, len(r.publications))
	copy(pubs, r.publications)
	return pubs, nil
}

func (r *publicationRepository) FindByID(ctx context.Context, id string) (*lcp.Publication, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, pub := range r.publications {
		if pub.ID == id {
			return pub, nil
		}
	}

	return nil, nil
}
