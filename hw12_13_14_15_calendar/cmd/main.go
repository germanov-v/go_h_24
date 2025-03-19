package main

import (
	"flag"
	"fmt"
	"github.com/germanov-v/go_h_24/hw12_13_14_15_calendar/configs"
	"github.com/germanov-v/go_h_24/hw12_13_14_15_calendar/pkg/middleware"
	"net/http"
)

// TODO: PANIC!
func App(pathConfig string) http.Handler {

	_, err := configs.LoadConfig(pathConfig)

	// https://golangify.com/routing-servemux
	router := http.NewServeMux()

	if err != nil {
		panic(err)
	}

	middlewareHandlers := middleware.FactoryMiddleware(
		middleware.LogProvider,
	)
	return middlewareHandlers(router)
}

func main() {

	path := flag.String("config", "config.json", "path to config file")
	flag.Parse()
	app := App(*path)

	server := http.Server{

		Addr:    ":8080",
		Handler: app,
	}
	fmt.Println("Start server")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
