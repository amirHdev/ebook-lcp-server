package lcp

import "time"

// Publication represents an encrypted book stored by the service.
type Publication struct {
	ID                  string    `db:"id" json:"id"`
	Title               string    `db:"title" json:"title"`
	Authors             []string  `db:"authors" json:"authors,omitempty"`
	Language            string    `db:"language" json:"language,omitempty"`
	Subjects            []string  `db:"subjects" json:"subjects,omitempty"`
	Tags                []string  `db:"tags" json:"tags,omitempty"`
	Status              string    `db:"status" json:"status,omitempty"`
	FilePath            string    `db:"file_path" json:"file_path,omitempty"`
	EncryptedPath       string    `db:"encrypted_path" json:"encrypted_path,omitempty"`
	EncryptedURI        string    `db:"encrypted_uri" json:"encrypted_uri,omitempty"`
	Checksum            string    `db:"checksum" json:"checksum,omitempty"`
	LicenseDurationDays int       `db:"license_duration_days" json:"licenseDurationDays,omitempty"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
}
