package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/vviveksharma/llmgate/internal/auth"
)

type ctxTenantKey struct{}

func TenantIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(ctxTenantKey{}).(string)
	return v
}

type Auth struct {
	APIKeys *auth.APIKeyValidator
	JWT     *auth.JWTValidator
}

func (a Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := ""
		if key := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "); key != r.Header.Get("Authorization") {
			if claims, err := a.JWT.Parse(strings.TrimSpace(key)); err == nil {
				tenantID = claims.TenantID
			}
		}
		if tenantID == "" {
			if tid, ok := a.APIKeys.Validate(r.Header.Get("X-API-Key")); ok {
				tenantID = tid
			}
		}
		if tenantID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ctxTenantKey{}, tenantID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
