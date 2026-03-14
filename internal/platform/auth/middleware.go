package auth

import (
	"effective_mobile_test/internal/platform/httputil"
	"net/http"
)

const APIKeyHeader = "Authorization"

func Middleware(expectedKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get(APIKeyHeader)
			if header == expectedKey {
				next.ServeHTTP(w, r)
				return
			}
			writeUnauthorized(w)
		})
	}
}

func writeUnauthorized(w http.ResponseWriter) {
	httputil.WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid or missing api key")
}
