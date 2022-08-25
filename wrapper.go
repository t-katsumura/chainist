package chainist

import "net/http"

type HandlerFuncWrapper struct {
	f http.HandlerFunc
}

func (h *HandlerFuncWrapper) PreMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.f != nil {
			h.f(w, r)
		}
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func (h *HandlerFuncWrapper) PostMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if next != nil {
			next.ServeHTTP(w, r)
		}
		if h.f != nil {
			h.f(w, r)
		}
	})
}

func (h *HandlerFuncWrapper) Middleware(next http.Handler) http.Handler {
	return h.PreMiddleware(next)
}
