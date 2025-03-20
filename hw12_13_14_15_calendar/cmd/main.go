package main

import (
	"flag"
	"fmt"
	"github.com/germanov-v/go_h_24/hw12_13_14_15_calendar/configs"
	"github.com/germanov-v/go_h_24/hw12_13_14_15_calendar/internal/home"
	"github.com/germanov-v/go_h_24/hw12_13_14_15_calendar/pkg/middleware"
	"log"
	"net/http"
	"os"
	"strconv"
)

// TODO: PANIC discuss!
func InitialApp(config *configs.Config) http.Handler {

	// https://golangify.com/routing-servemux
	router := http.NewServeMux()

	home.AddHomeBaseHandlers(router)

	middlewareHandlers := middleware.FactoryMiddleware(
		middleware.LogProvider,
	)
	return middlewareHandlers(router)
}

func main() {

	path := flag.String("config", "config.json", "path to config file")
	flag.Parse()
	config, err := configs.LoadConfig(*path)
	if err != nil {
		panic(err)
	}

	var logFile *os.File
	if config.LogConfig.Providers.File != "" {
		logFile, err = os.OpenFile(config.LogConfig.Providers.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("Error open file log occurred", err)
		}
		// здесь нужно замые
		defer logFile.Close()
		log.SetOutput(logFile)

	}

	app := InitialApp(config)

	server := http.Server{

		Addr:    config.ServerConfig.Host + ":" + strconv.Itoa(config.ServerConfig.Port),
		Handler: app,
	}

	fmt.Println("Start server")
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
