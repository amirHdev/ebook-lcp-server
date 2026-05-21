package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type config struct {
	baseURL   string
	username  string
	password  string
	twoFactor string
	client    *http.Client
}

type loginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

type publicationResult struct {
	ID          string `json:"id"`
	Title       string `json:"title,omitempty"`
	DownloadURL string `json:"downloadURL,omitempty"`
}

type licenseResult struct {
	ID             string `json:"id"`
	PublicationID  string `json:"publicationID"`
	UserID         string `json:"userID"`
	PublicationURL string `json:"publicationURL,omitempty"`
	Passphrase     string `json:"passphrase,omitempty"`
	Hint           string `json:"hint,omitempty"`
}

type graphQLResponse struct {
	Data   map[string]json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		printUsage(stderr)
		return errors.New("command is required")
	}

	cfg := newConfig()

	switch args[0] {
	case "health":
		return runHealth(cfg, args[1:], stdout)
	case "login":
		return runLogin(cfg, args[1:], stdout)
	case "upload":
		return runUpload(cfg, args[1:], stdout)
	case "license":
		return runLicense(cfg, args[1:], stdout, stderr)
	case "demo":
		return runDemo(cfg, args[1:], stdout)
	case "help", "-h", "--help":
		printUsage(stdout)
		return nil
	default:
		printUsage(stderr)
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func newConfig() config {
	baseURL := strings.TrimRight(envOr("LCP_BASE_URL", "http://127.0.0.1:8080"), "/")
	return config{
		baseURL:   baseURL,
		username:  envOr("LCP_USERNAME", "publisher"),
		password:  envOr("LCP_PASSWORD", "publisher"),
		twoFactor: envOr("LCP_2FA_CODE", ""),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func envOr(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage: lcpctl <command> [options]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  health                Check /healthz and /readyz")
	fmt.Fprintln(w, "  login                 Get a JWT for API access")
	fmt.Fprintln(w, "  upload                Upload a publication file")
	fmt.Fprintln(w, "  license create        Create a license for a publication")
	fmt.Fprintln(w, "  license revoke        Revoke an existing license")
	fmt.Fprintln(w, "  demo                  Upload a book and issue a demo license")
}

func addAuthFlags(fs *flag.FlagSet, cfg *config) {
	fs.StringVar(&cfg.baseURL, "base-url", cfg.baseURL, "LCP API base URL")
	fs.StringVar(&cfg.username, "username", cfg.username, "API username")
	fs.StringVar(&cfg.password, "password", cfg.password, "API password")
	fs.StringVar(&cfg.twoFactor, "two-factor", cfg.twoFactor, "2FA code for admin login")
}

func runHealth(cfg config, args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("health", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&cfg.baseURL, "base-url", cfg.baseURL, "LCP API base URL")
	if err := fs.Parse(args); err != nil {
		return err
	}
	ctx := context.Background()
	health, err := getJSON(ctx, cfg.client, cfg.baseURL+"/healthz", nil)
	if err != nil {
		return err
	}
	ready, err := getJSON(ctx, cfg.client, cfg.baseURL+"/readyz", nil)
	if err != nil {
		return err
	}
	return writeJSON(stdout, map[string]any{
		"baseURL": cfg.baseURL,
		"health":  health,
		"ready":   ready,
	})
}

func runLogin(cfg config, args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("login", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	tokenOnly := fs.Bool("token-only", false, "print only the JWT")
	addAuthFlags(fs, &cfg)
	if err := fs.Parse(args); err != nil {
		return err
	}
	login, err := authenticate(context.Background(), cfg)
	if err != nil {
		return err
	}
	if *tokenOnly {
		_, err = fmt.Fprintln(stdout, login.Token)
		return err
	}
	return writeJSON(stdout, login)
}

func runUpload(cfg config, args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("upload", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	filePath := fs.String("file", "", "Path to EPUB or PDF file")
	title := fs.String("title", "", "Publication title")
	addAuthFlags(fs, &cfg)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*filePath) == "" {
		return errors.New("--file is required")
	}
	publication, err := uploadPublication(context.Background(), cfg, *filePath, *title)
	if err != nil {
		return err
	}
	return writeJSON(stdout, publication)
}

func runLicense(cfg config, args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "license command requires create or revoke")
		return errors.New("license subcommand is required")
	}
	switch args[0] {
	case "create":
		return runCreateLicense(cfg, args[1:], stdout)
	case "revoke":
		return runRevokeLicense(cfg, args[1:], stdout)
	default:
		return fmt.Errorf("unknown license subcommand: %s", args[0])
	}
}

func runCreateLicense(cfg config, args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("license create", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	publicationID := fs.String("publication-id", "", "Publication ID")
	userID := fs.String("user-id", "reader-01", "Reader identifier")
	passphrase := fs.String("passphrase", "open-sesame", "License passphrase")
	hint := fs.String("hint", "demo", "Passphrase hint")
	addAuthFlags(fs, &cfg)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*publicationID) == "" {
		return errors.New("--publication-id is required")
	}
	license, err := createLicense(context.Background(), cfg, *publicationID, *userID, *passphrase, *hint)
	if err != nil {
		return err
	}
	return writeJSON(stdout, license)
}

func runRevokeLicense(cfg config, args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("license revoke", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	licenseID := fs.String("id", "", "License ID")
	addAuthFlags(fs, &cfg)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*licenseID) == "" {
		return errors.New("--id is required")
	}
	if err := revokeLicense(context.Background(), cfg, *licenseID); err != nil {
		return err
	}
	return writeJSON(stdout, map[string]string{
		"status":    "revoked",
		"licenseID": *licenseID,
	})
}

func runDemo(cfg config, args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("demo", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	filePath := fs.String("file", "examples/pride-and-prejudice/pride-and-prejudice.epub", "Path to EPUB or PDF file")
	title := fs.String("title", "Pride and Prejudice", "Publication title")
	userID := fs.String("user-id", "reader-01", "Reader identifier")
	passphrase := fs.String("passphrase", "open-sesame", "License passphrase")
	hint := fs.String("hint", "demo", "Passphrase hint")
	addAuthFlags(fs, &cfg)
	if err := fs.Parse(args); err != nil {
		return err
	}
	publication, err := uploadPublication(context.Background(), cfg, *filePath, *title)
	if err != nil {
		return err
	}
	license, err := createLicense(context.Background(), cfg, publication.ID, *userID, *passphrase, *hint)
	if err != nil {
		return err
	}
	return writeJSON(stdout, map[string]any{
		"publication_id": publication.ID,
		"license_id":     license.ID,
		"license_url":    strings.TrimRight(cfg.baseURL, "/") + "/licenses/" + license.ID + ".lcpl",
	})
}

func authenticate(ctx context.Context, cfg config) (*loginResponse, error) {
	payload := map[string]string{
		"username":  cfg.username,
		"password":  cfg.password,
		"twoFactor": cfg.twoFactor,
	}
	body, err := postJSON(ctx, cfg.client, cfg.baseURL+"/api/v1/auth/login", payload, "")
	if err != nil {
		return nil, err
	}
	var response loginResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	if response.Token == "" {
		return nil, errors.New("empty token in login response")
	}
	return &response, nil
}

func uploadPublication(ctx context.Context, cfg config, filePath, title string) (*publicationResult, error) {
	raw, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(title) == "" {
		title = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	}
	login, err := authenticate(ctx, cfg)
	if err != nil {
		return nil, err
	}
	payload := map[string]any{
		"query": "mutation UploadPublication($title: String!, $file: Upload!) { uploadPublication(title: $title, file: $file) { id title downloadURL } }",
		"variables": map[string]string{
			"title": title,
			"file":  base64.StdEncoding.EncodeToString(raw),
		},
	}
	body, err := postJSON(ctx, cfg.client, cfg.baseURL+"/graphql", payload, login.Token)
	if err != nil {
		return nil, err
	}
	var response graphQLResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	if err := firstGraphQLError(response); err != nil {
		return nil, err
	}
	var publication publicationResult
	if err := json.Unmarshal(response.Data["uploadPublication"], &publication); err != nil {
		return nil, err
	}
	return &publication, nil
}

func createLicense(ctx context.Context, cfg config, publicationID, userID, passphrase, hint string) (*licenseResult, error) {
	login, err := authenticate(ctx, cfg)
	if err != nil {
		return nil, err
	}
	payload := map[string]any{
		"query": "mutation CreateLicense($publicationID: ID!, $userID: ID!, $passphrase: String!, $hint: String!) { createLicense(publicationID: $publicationID, userID: $userID, passphrase: $passphrase, hint: $hint) { id publicationID userID publicationURL passphrase hint } }",
		"variables": map[string]string{
			"publicationID": publicationID,
			"userID":        userID,
			"passphrase":    passphrase,
			"hint":          hint,
		},
	}
	body, err := postJSON(ctx, cfg.client, cfg.baseURL+"/graphql", payload, login.Token)
	if err != nil {
		return nil, err
	}
	var response graphQLResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	if err := firstGraphQLError(response); err != nil {
		return nil, err
	}
	var license licenseResult
	if err := json.Unmarshal(response.Data["createLicense"], &license); err != nil {
		return nil, err
	}
	return &license, nil
}

func revokeLicense(ctx context.Context, cfg config, licenseID string) error {
	login, err := authenticate(ctx, cfg)
	if err != nil {
		return err
	}
	payload := map[string]any{
		"query": "mutation RevokeLicense($id: ID!) { revokeLicense(id: $id) }",
		"variables": map[string]string{
			"id": licenseID,
		},
	}
	body, err := postJSON(ctx, cfg.client, cfg.baseURL+"/graphql", payload, login.Token)
	if err != nil {
		return err
	}
	var response graphQLResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}
	if err := firstGraphQLError(response); err != nil {
		return err
	}
	if raw, ok := response.Data["revokeLicense"]; ok && string(raw) == "true" {
		return nil
	}
	return nil
}

func firstGraphQLError(response graphQLResponse) error {
	if len(response.Errors) == 0 {
		return nil
	}
	return errors.New(response.Errors[0].Message)
}

func postJSON(ctx context.Context, client *http.Client, url string, payload any, token string) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s returned %s: %s", url, resp.Status, strings.TrimSpace(string(body)))
	}
	return body, nil
}

func getJSON(ctx context.Context, client *http.Client, url string, headers map[string]string) (map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s returned %s: %s", url, resp.Status, strings.TrimSpace(string(body)))
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func writeJSON(w io.Writer, payload any) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(payload)
}
