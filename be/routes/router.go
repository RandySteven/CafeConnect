package routes

import (
	"github.com/RandySteven/CafeConnect/be/enums"
	"net/http"
)

type (
	HandlerFunc func(w http.ResponseWriter, r *http.Request)

	Router struct {
		path        string
		handler     HandlerFunc
		method      string
		middlewares []enums.Middleware
	}

	RouterPrefix map[enums.RouterPrefix][]*Router
)
