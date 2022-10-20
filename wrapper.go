package chainist

import "net/http"

type HandlerFuncWrapper struct {
	HandlerFunc http.HandlerFunc
}

// Wrap http handler function as http handler.
// Wrapped function is executed before invoking proceeding handlers.
// This is what `Pre` means.
func (h *HandlerFuncWrapper) PreMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.HandlerFunc != nil {
			h.HandlerFunc(w, r)
		}
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func GetPreMiddleware(handlerFunc http.HandlerFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if handlerFunc != nil {
				handlerFunc(w, r)
			}
			if next != nil {
				next.ServeHTTP(w, r)
			}
		})
	}
}

// Wrap http handler function as http handler.
// Wrapped function is executed after invoking proceeding handlers.
// This is what `Post` means.
func (h *HandlerFuncWrapper) PostMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if next != nil {
			next.ServeHTTP(w, r)
		}
		if h.HandlerFunc != nil {
			h.HandlerFunc(w, r)
		}
	})
}

func GetPostMiddleware(handlerFunc http.HandlerFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if next != nil {
				next.ServeHTTP(w, r)
			}
			if handlerFunc != nil {
				handlerFunc(w, r)
			}
		})
	}
}

// Alias for `PreMiddleware`
func (h *HandlerFuncWrapper) Middleware(next http.Handler) http.Handler {
	return h.PreMiddleware(next)
}
