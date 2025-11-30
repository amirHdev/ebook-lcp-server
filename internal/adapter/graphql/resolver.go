package graphql

import (
	usecaseLicense "github.com/Mehrbod2002/lcp/internal/usecase/lcp/license"
	usecasePublication "github.com/Mehrbod2002/lcp/internal/usecase/lcp/publication"
)

// Resolver aggregates the use cases needed by the GraphQL layer.
type Resolver struct {
	PublicationUsecase usecasePublication.PublicationUsecase
	LicenseUsecase     usecaseLicense.LicenseUsecase
	PublicBaseURL      string
}
