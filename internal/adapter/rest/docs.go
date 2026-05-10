package rest

import (
	"net/http"
	"path/filepath"
)

func SwaggerYAML() http.HandlerFunc {
	return staticFileHandler("docs/swagger.yaml", "text/yaml; charset=utf-8")
}

func SwaggerJSON() http.HandlerFunc {
	return staticFileHandler("docs/swagger.json", "application/json; charset=utf-8")
}

func OpenAPIYAML() http.HandlerFunc {
	return staticFileHandler("docs/openapi-rest.yaml", "text/yaml; charset=utf-8")
}

func staticFileHandler(path string, contentType string) http.HandlerFunc {
	normalized := filepath.Clean(path)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		w.Header().Set("Content-Type", contentType)
		http.ServeFile(w, r, normalized)
	}
}
