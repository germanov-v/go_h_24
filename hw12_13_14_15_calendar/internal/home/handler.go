package home

import (
	"github.com/germanov-v/go_h_24/hw12_13_14_15_calendar/pkg/storage_v2"
	"net/http"
	"time"
)

type dataHandler struct {
	storage storage_v2.Storage
	format  string
}

func AddHomeBaseHandlers(router *http.ServeMux,

	storage storage_v2.Storage) {

	th := dataHandler{storage: storage, format: time.RFC3339}
	router.Handle("GET /", th)
}

type timeHandler struct {
	format string
}

// https://www.alexedwards.net/blog/an-introduction-to-handlers-and-servemuxes-in-go
func (dH dataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	//event := storage_v2.EventItem{}
	//_, err := dH.storage.CreateEvent(event, ctx)
	//if err != nil {
	//	return
	//}
	tm := time.Now().Format(dH.format)
	w.Write([]byte("The time is: " + tm))
}
