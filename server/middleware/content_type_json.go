package middleware

import "net/http"

// ContentTypeJSON simply sets the header of response to application/json
func ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json;charset=utf8")
		next.ServeHTTP(rw, r)
	})
}
