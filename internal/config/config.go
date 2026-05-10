package config

import (
	"os"
)

type Config struct {
	Database struct {
		DSN string
	}
	LCP struct {
		Profile     string // "basic" or "production"
		Certificate string // Path to X.509 certificate
		PrivateKey  string // Path to private key
		ProviderURI string
		CoreURL     string
		CoreUser    string
		CorePass    string
		Storage     struct {
			Mode string // "fs" or "s3"
			FS   struct {
				Directory string
			}
			S3 struct {
				Region    string
				Bucket    string
				AccessKey string
				SecretKey string
			}
		}
	}
	JWT struct {
		Secret       string
		Admin2FACode string
	}
	Admin struct {
		Username string
		Password string
	}
	Server struct {
		Port          string
		PublicBaseURL string
		StatusBaseURL string
	}
	DataDir string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	cfg.Database.DSN = os.Getenv("DB_DSN")
	cfg.LCP.Profile = os.Getenv("LCP_PROFILE")
	cfg.LCP.Certificate = os.Getenv("LCP_CERTIFICATE")
	cfg.LCP.PrivateKey = os.Getenv("LCP_PRIVATE_KEY")
	cfg.LCP.ProviderURI = os.Getenv("LCP_PROVIDER_URI")
	cfg.LCP.CoreURL = os.Getenv("LCP_CORE_URL")
	cfg.LCP.CoreUser = os.Getenv("LCP_CORE_USER")
	cfg.LCP.CorePass = os.Getenv("LCP_CORE_PASSWORD")
	cfg.LCP.Storage.Mode = os.Getenv("LCP_STORAGE_MODE")
	cfg.LCP.Storage.FS.Directory = os.Getenv("LCP_STORAGE_FS_DIR")
	cfg.LCP.Storage.S3.Region = os.Getenv("LCP_S3_REGION")
	cfg.LCP.Storage.S3.Bucket = os.Getenv("LCP_S3_BUCKET")
	cfg.LCP.Storage.S3.AccessKey = os.Getenv("LCP_S3_ACCESS_KEY")
	cfg.LCP.Storage.S3.SecretKey = os.Getenv("LCP_S3_SECRET_KEY")
	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	cfg.JWT.Admin2FACode = os.Getenv("ADMIN_2FA_CODE")
	cfg.Admin.Username = os.Getenv("ADMIN_USERNAME")
	cfg.Admin.Password = os.Getenv("ADMIN_PASSWORD")
	cfg.Server.Port = os.Getenv("SERVER_PORT")
	cfg.Server.PublicBaseURL = os.Getenv("PUBLIC_BASE_URL")
	cfg.Server.StatusBaseURL = os.Getenv("STATUS_BASE_URL")
	cfg.DataDir = os.Getenv("DATA_DIR")
	return cfg, nil
}
