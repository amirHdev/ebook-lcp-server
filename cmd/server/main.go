package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/Mehrbod2002/lcp/internal/adapter/graphql"
	"github.com/Mehrbod2002/lcp/internal/adapter/repository/lcp"
	"github.com/Mehrbod2002/lcp/internal/config"
	lcpencrypt "github.com/Mehrbod2002/lcp/internal/lcp/encrypt"
	lcplicense "github.com/Mehrbod2002/lcp/internal/lcp/license"
	"github.com/Mehrbod2002/lcp/internal/usecase/lcp/license"
	"github.com/Mehrbod2002/lcp/internal/usecase/lcp/publication"
)

// @title LCP License Server API
// @version 1.0
// @description API for managing LCP licenses and publications
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	lcpEnc := lcpencrypt.NewFileCopyEncrypter(cfg.LCP.Storage.FS.Directory)
	lcpSrv := lcplicense.NewService()
	pubRepo := lcp.NewPublicationRepository()
	licRepo := lcp.NewLicenseRepository()
	pubUsecase := publication.NewPublicationUsecase(pubRepo, lcpEnc)
	publicBaseURL := buildBaseURL(cfg)
	licUsecase := license.NewLicenseUsecase(licRepo, lcpSrv, publicBaseURL)

	mux := http.NewServeMux()

	gqlHandler := graphql.NewHandler(&graphql.Resolver{
		PublicationUsecase: pubUsecase,
		LicenseUsecase:     licUsecase,
		PublicBaseURL:      publicBaseURL,
	})
	mux.Handle("/graphql", gqlHandler)
	mux.Handle("/publications/", publicationDownloadHandler(pubUsecase))

	port := cfg.Server.Port
	if port == "" {
		port = ":8080"
	}

	if err := http.ListenAndServe(port, mux); err != nil {
		panic(err)
	}
}

func buildBaseURL(cfg *config.Config) string {
	baseURL := strings.TrimSpace(cfg.Server.PublicBaseURL)
	if baseURL != "" {
		return strings.TrimSuffix(baseURL, "/")
	}

	port := cfg.Server.Port
	if port == "" {
		port = ":8080"
	}

	return "http://localhost" + port
}

func publicationDownloadHandler(pubUsecase publication.PublicationUsecase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 3 || parts[0] != "publications" || parts[2] != "content" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		pubID := parts[1]
		pub, err := pubUsecase.GetByID(context.Background(), pubID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if pub == nil || pub.EncryptedPath == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		http.ServeFile(w, r, pub.EncryptedPath)
	})
}
