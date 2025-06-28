package middlewares

import "net/http"

func (s *ServerMiddleware) IdempotencyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Get(`idempotency-key`)

		next.ServeHTTP(w, r)
	})
}
