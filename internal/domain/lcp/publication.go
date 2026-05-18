package lcp

import "time"

// Publication represents an encrypted book stored by the service.
type Publication struct {
	ID                  string    `db:"id" json:"id"`
	TenantID            string    `db:"tenant_id" json:"tenantId,omitempty"`
	Title               string    `db:"title" json:"title"`
	Authors             []string  `db:"authors" json:"authors,omitempty"`
	Language            string    `db:"language" json:"language,omitempty"`
	Subjects            []string  `db:"subjects" json:"subjects,omitempty"`
	Tags                []string  `db:"tags" json:"tags,omitempty"`
	Status              string    `db:"status" json:"status,omitempty"`
	RightPrint          *int      `db:"right_print" json:"right_print,omitempty"`
	RightCopy           *int      `db:"right_copy" json:"right_copy,omitempty"`
	FilePath            string    `db:"file_path" json:"file_path,omitempty"`
	EncryptedPath       string    `db:"encrypted_path" json:"encrypted_path,omitempty"`
	EncryptedURI        string    `db:"encrypted_uri" json:"encrypted_uri,omitempty"`
	Checksum            string    `db:"checksum" json:"checksum,omitempty"`
	LicenseDurationDays int       `db:"license_duration_days" json:"licenseDurationDays,omitempty"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
}
