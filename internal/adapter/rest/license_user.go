package rest

import (
	"net/http"
	"strings"

	usecaseLicense "github.com/Mehrbod2002/lcp/internal/usecase/lcp/license"
)

type userDataPayload struct {
	ID             string `json:"id"`
	Name           string `json:"name,omitempty"`
	Email          string `json:"email,omitempty"`
	PassphraseHash string `json:"passphrasehash"`
	Hint           string `json:"hint,omitempty"`
}

func LicenseUserData(licenses usecaseLicense.LicenseUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		licenseID := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/licenses/"), "/user")
		licenseID = strings.Trim(licenseID, "/")
		if licenseID == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "license id is required"})
			return
		}

		lic, err := licenses.GetByID(r.Context(), licenseID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if lic == nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "license not found"})
			return
		}

		payload := userDataPayload{
			ID:             lic.UserID,
			PassphraseHash: lic.PassphraseHash,
			Hint:           lic.Hint,
		}

		writeJSON(w, http.StatusOK, payload)
	}
}
