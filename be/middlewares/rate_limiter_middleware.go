package middlewares

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/enums"
	redis_client "github.com/RandySteven/CafeConnect/be/pkg/redis"
	"github.com/RandySteven/CafeConnect/be/utils"
	ip "github.com/vikram1565/request-ip"
	"net/http"
)

func (c *ClientMiddleware) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !c.middlewares.WhiteListed(r.Method, utils.ReplaceLastURLID(r.RequestURI), enums.RateLimiterMiddleware) {
			next.ServeHTTP(w, r)
			return
		}
		clientIp := ip.GetClientIP(r)
		ctx := context.WithValue(r.Context(), enums.ClientIP, clientIp)
		if err := redis_client.RateLimiter(ctx); err != nil {
			w.WriteHeader(http.StatusTooManyRequests)
			utils.ResponseHandler(w, http.StatusTooManyRequests, `too many request`, nil, nil, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}
