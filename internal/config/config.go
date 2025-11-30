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
		Secret string
	}
	Server struct {
		Port          string
		PublicBaseURL string
	}
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	cfg.Database.DSN = os.Getenv("DB_DSN")
	cfg.LCP.Profile = os.Getenv("LCP_PROFILE")
	cfg.LCP.Certificate = os.Getenv("LCP_CERTIFICATE")
	cfg.LCP.PrivateKey = os.Getenv("LCP_PRIVATE_KEY")
	cfg.LCP.Storage.Mode = os.Getenv("LCP_STORAGE_MODE")
	cfg.LCP.Storage.FS.Directory = os.Getenv("LCP_STORAGE_FS_DIR")
	cfg.LCP.Storage.S3.Region = os.Getenv("LCP_S3_REGION")
	cfg.LCP.Storage.S3.Bucket = os.Getenv("LCP_S3_BUCKET")
	cfg.LCP.Storage.S3.AccessKey = os.Getenv("LCP_S3_ACCESS_KEY")
	cfg.LCP.Storage.S3.SecretKey = os.Getenv("LCP_S3_SECRET_KEY")
	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	cfg.Server.Port = os.Getenv("SERVER_PORT")
	cfg.Server.PublicBaseURL = os.Getenv("PUBLIC_BASE_URL")
	return cfg, nil
}
