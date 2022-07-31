package middleware

import (
	"net/http"
	"strings"
)

// RemoveTrailingSlash - because mux does not treat /foo the same as /foo/
// so you may get annoying 405s. Solution taken from:
// https://www.husainalshehhi.com/blog/gorilla-mux-trailing-slashes/
func RemoveTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		next.ServeHTTP(w, r)
	})
}
