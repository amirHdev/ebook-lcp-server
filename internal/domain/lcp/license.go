package lcp

import "time"

// License captures access information for a publication.
type License struct {
	ID             string     `db:"id" json:"id"`
	PublicationID  string     `db:"publication_id" json:"publication_id"`
	UserID         string     `db:"user_id" json:"user_id"`
	Passphrase     string     `db:"passphrase" json:"passphrase"`
	Hint           string     `db:"hint" json:"hint"`
	PublicationURL string     `db:"publication_url" json:"publication_url"`
	RightPrint     *int       `db:"right_print" json:"right_print"`
	RightCopy      *int       `db:"right_copy" json:"right_copy"`
	StartDate      *time.Time `db:"start_date" json:"start_date"`
	EndDate        *time.Time `db:"end_date" json:"end_date"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
}

// LicenseInput is the input contract for creating a license.
type LicenseInput struct {
	PublicationID string     `json:"publication_id"`
	UserID        string     `json:"user_id"`
	Passphrase    string     `json:"passphrase"`
	Hint          string     `json:"hint"`
	RightPrint    *int       `json:"right_print"`
	RightCopy     *int       `json:"right_copy"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
}
