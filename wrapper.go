package chainist

import "net/http"

type HandlerFuncWrapper struct {
	f http.HandlerFunc
}

func (h *HandlerFuncWrapper) preMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.f != nil {
			h.f(w, r)
		}
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func (h *HandlerFuncWrapper) postMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if next != nil {
			next.ServeHTTP(w, r)
		}
		if h.f != nil {
			h.f(w, r)
		}
	})
}

func (h *HandlerFuncWrapper) middleware(next http.Handler) http.Handler {
	return h.preMiddleware(next)
}
