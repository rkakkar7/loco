package main

import (
	"fmt"
	"loco/router"
	"log"
	"net/http"
	"time"
)

func main() {
	r := router.Router()
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Starting server on the port 8080...")
	log.Fatal(srv.ListenAndServe())
}
