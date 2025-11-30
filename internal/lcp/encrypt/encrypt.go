package encrypt

import (
	"io"
	"os"
	"path/filepath"
)

// Encrypter defines the behavior required by the publication use case.
type Encrypter interface {
	Encrypt(inputPath, label string) (string, error)
}

// FileCopyEncrypter is a simple stand-in that copies the source file to a
// destination path to simulate encryption. This keeps the rest of the system
// testable without binding to a specific DRM backend.
type FileCopyEncrypter struct {
	OutputDir string
}

// Encrypt copies the input file into the configured output directory and
// returns the resulting path.
func (e *FileCopyEncrypter) Encrypt(inputPath, label string) (string, error) {
	if err := os.MkdirAll(e.OutputDir, 0o755); err != nil {
		return "", err
	}

	dest := filepath.Join(e.OutputDir, filepath.Base(inputPath))
	in, err := os.Open(inputPath)
	if err != nil {
		return "", err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return "", err
	}

	return dest, nil
}

// NewFileCopyEncrypter constructs the default encrypter used in development
// and testing environments.
func NewFileCopyEncrypter(outputDir string) *FileCopyEncrypter {
	return &FileCopyEncrypter{OutputDir: outputDir}
}
