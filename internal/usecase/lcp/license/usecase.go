package license

import (
	"context"
	"time"

	"github.com/Mehrbod2002/lcp/internal/domain/lcp"
	lcplicense "github.com/Mehrbod2002/lcp/internal/lcp/license"
	"github.com/Mehrbod2002/lcp/internal/pkg/id"
)

type LicenseUsecase interface {
	Create(ctx context.Context, input *lcp.LicenseInput) (*lcp.License, error)
	GetByPublication(ctx context.Context, publicationID *string) ([]*lcp.License, error)
	Revoke(ctx context.Context, id string) error
}

type licenseUsecase struct {
	repo    lcp.LicenseRepository
	lcp     *lcplicense.Service
	baseURL string
}

func NewLicenseUsecase(repo lcp.LicenseRepository, lcp *lcplicense.Service, baseURL string) LicenseUsecase {
	return &licenseUsecase{repo: repo, lcp: lcp, baseURL: baseURL}
}

func (u *licenseUsecase) Create(ctx context.Context, input *lcp.LicenseInput) (*lcp.License, error) {
	license := &lcp.License{
		ID:             id.New(),
		PublicationID:  input.PublicationID,
		UserID:         input.UserID,
		Passphrase:     input.Passphrase,
		Hint:           input.Hint,
		PublicationURL: u.baseURL + "/publications/" + input.PublicationID + "/content",
		RightPrint:     input.RightPrint,
		RightCopy:      input.RightCopy,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		CreatedAt:      time.Now(),
	}

	// Generate LCP license using lcpserver
	err := u.lcp.GenerateLicense(license)
	if err != nil {
		return nil, err
	}

	// Save license to database
	err = u.repo.Save(ctx, license)
	if err != nil {
		return nil, err
	}

	return license, nil
}

func (u *licenseUsecase) GetByPublication(ctx context.Context, publicationID *string) ([]*lcp.License, error) {
	return u.repo.FindByPublication(ctx, publicationID)
}

func (u *licenseUsecase) Revoke(ctx context.Context, id string) error {
	return u.lcp.RevokeLicense(id)
}
