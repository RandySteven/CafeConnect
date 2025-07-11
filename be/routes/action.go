package routes

import (
	"github.com/RandySteven/CafeConnect/be/enums"
	"net/http"
)

func registerEndpointRouter(path, method string, handler HandlerFunc, middlewares ...enums.Middleware) *Router {
	return &Router{path: path, handler: handler, method: method, middlewares: middlewares}
}

func Post(path string, handler HandlerFunc, middlewares ...enums.Middleware) *Router {
	return registerEndpointRouter(path, http.MethodPost, handler, middlewares...)
}

func Get(path string, handler HandlerFunc, middlewares ...enums.Middleware) *Router {
	return registerEndpointRouter(path, http.MethodGet, handler, middlewares...)
}

func Put(path string, handler HandlerFunc, middlewares ...enums.Middleware) *Router {
	return registerEndpointRouter(path, http.MethodPut, handler, middlewares...)
}

func Delete(path string, handler HandlerFunc, middlewares ...enums.Middleware) *Router {
	return registerEndpointRouter(path, http.MethodDelete, handler, middlewares...)
}
