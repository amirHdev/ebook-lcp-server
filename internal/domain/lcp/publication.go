package lcp

import "time"

// Publication represents an encrypted book stored by the service.
type Publication struct {
	ID            string    `db:"id" json:"id"`
	Title         string    `db:"title" json:"title"`
	FilePath      string    `db:"file_path" json:"file_path"`
	EncryptedPath string    `db:"encrypted_path" json:"encrypted_path"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}
