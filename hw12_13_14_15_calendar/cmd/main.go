package main

import (
	"flag"
	"fmt"
	"github.com/germanov-v/go_h_24/hw12_13_14_15_calendar/configs"
	"github.com/germanov-v/go_h_24/hw12_13_14_15_calendar/internal/home"
	"github.com/germanov-v/go_h_24/hw12_13_14_15_calendar/pkg/middleware"
	"github.com/germanov-v/go_h_24/hw12_13_14_15_calendar/pkg/storage_v2"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"

	"log"
	"net/http"
	"os"
	"strconv"
)

// TODO: PANIC discuss!
func InitialApp(config *configs.Config) http.Handler {

	// https://golangify.com/routing-servemux
	router := http.NewServeMux()

	var err error
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

	var store storage_v2.Storage

	if config.DbBaseConfig.ConnectionString != "" {
		db, err := sqlx.Connect("postgres", config.DbBaseConfig.ConnectionString)
		if err != nil {
			log.Fatal("Error connecting to database", err)
		}
		store = storage_v2.NewPostgresStorage(db)
	} else {
		store = storage_v2.NewMemoryStorage()
	}

	home.AddHomeBaseHandlers(router, store)

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
