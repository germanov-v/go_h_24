package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

type ResponseWrapper struct {
	StatusCode int
	http.ResponseWriter
}

func FactoryMiddleware(middlewares ...Middleware) Middleware {
	f := func(next http.Handler) http.Handler {
		//for _, m := range middlewares {
		//	next = m(next)
		//}
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
	return f
}
