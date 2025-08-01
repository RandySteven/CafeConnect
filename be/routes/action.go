package routes

import (
	"github.com/RandySteven/CafeConnect/be/enums"
	"net/http"
)

func registerEndpointRouter(methodName string, path, method string, handler HandlerFunc, middlewares ...enums.Middleware) *Router {
	return &Router{methodName: methodName, path: path, handler: handler, method: method, middlewares: middlewares}
}

func Post(methodName string, path string, handler HandlerFunc, middlewares ...enums.Middleware) *Router {
	return registerEndpointRouter(methodName, path, http.MethodPost, handler, middlewares...)
}

func Get(methodName string, path string, handler HandlerFunc, middlewares ...enums.Middleware) *Router {
	return registerEndpointRouter(methodName, path, http.MethodGet, handler, middlewares...)
}

func Put(methodName string, path string, handler HandlerFunc, middlewares ...enums.Middleware) *Router {
	return registerEndpointRouter(methodName, path, http.MethodPut, handler, middlewares...)
}

func Delete(methodName string, path string, handler HandlerFunc, middlewares ...enums.Middleware) *Router {
	return registerEndpointRouter(methodName, path, http.MethodDelete, handler, middlewares...)
}
