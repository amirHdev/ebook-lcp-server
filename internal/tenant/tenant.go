package tenant

import (
	"context"

	"github.com/amirhdev/ebook-lcp-server/internal/auth"
)

func IDFromContext(ctx context.Context) string {
	if claims, ok := auth.FromContext(ctx); ok && claims.TenantID != "" {
		return claims.TenantID
	}
	return "default"
}
