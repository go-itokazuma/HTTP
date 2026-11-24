package main

import (
	"net/http"
	"strings"
)

func methodOverrideMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err == nil {
				override := strings.ToUpper((r.Form.Get("_method")))
				if override == http.MethodPut || override == http.MethodDelete {
					r.Method = override
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
