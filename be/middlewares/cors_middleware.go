package middlewares

import (
	"github.com/RandySteven/CafeConnect/be/enums"
	"github.com/RandySteven/CafeConnect/be/utils"
	"net/http"
)

func (s *ServerMiddleware) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Origin, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		}
		utils.ContentType(w, enums.ContentTypeJSON)
		next.ServeHTTP(w, r)
	})
}
