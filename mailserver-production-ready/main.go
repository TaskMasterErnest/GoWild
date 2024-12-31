package main

import (
	"log"
	"time"

	smtp "github.com/emersion/go-smtp"
)

func main() {
	server := smtp.NewServer(&Backend{})
	// add server attributes
	server.Addr = ":2525"
	server.Domain = "localhost" // change to prod domain
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.MaxMessageBytes = 1024 * 1024
	server.MaxRecipients = 50
	server.AllowInsecureAuth = false

	// start the server
	log.Println("Starting Mail Server at ", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
