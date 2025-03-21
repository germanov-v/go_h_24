package middleware

import (
	"fmt"
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

		//log.Printf("%s - %s - %s - %s", now, r.Method, r.URL.Path, wrapper.StatusCode)
		//log.Println(wrapper.StatusCode, now, now.Sub(now))
		ip := r.RemoteAddr
		logEntry := fmt.Sprintf(
			`%s [%s] "%s %s %s" %d %s "%s"`,
			ip,
			now.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.URL.RequestURI(),
			r.Proto,
			wrapper.StatusCode,
			w.Header().Get("Content-Length"),
			r.UserAgent(),
		)

		log.Println(logEntry)

	})
	return handler
}
