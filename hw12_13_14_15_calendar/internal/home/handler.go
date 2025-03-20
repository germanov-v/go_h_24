package home

import (
	"net/http"
	"time"
)

func AddHomeBaseHandlers(router *http.ServeMux) {

	th := timeHandler{format: time.RFC3339}
	router.Handle("GET /", th)
}

type timeHandler struct {
	format string
}

// https://www.alexedwards.net/blog/an-introduction-to-handlers-and-servemuxes-in-go
func (th timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(th.format)
	w.Write([]byte("The time is: " + tm))
}
