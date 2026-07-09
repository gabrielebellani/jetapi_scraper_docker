package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger

	apiCalls     atomic.Uint64
	totalLatency atomic.Int64 // stored as nanoseconds
}

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8922"
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	app.infoLog.Print("Starting stats logger")
	go app.statsLogger()

	app.infoLog.Printf("Starting server on %s", addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
