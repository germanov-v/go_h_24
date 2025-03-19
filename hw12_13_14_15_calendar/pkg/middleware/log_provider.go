package middleware

import (
	"log"
	"net/http"
	"time"
)

func LogProvider(next http.Handler) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().UTC()
		wrapper := &ResponseWrapper{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		// in
		next.ServeHTTP(wrapper, r)
		// out

		log.Printf("%s - %s - %s - %s", now, r.Method, r.URL.Path, wrapper.StatusCode)
		log.Println(wrapper.StatusCode, now, now.Sub(now))

	})
	return handler
}
