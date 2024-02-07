package main

import (
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
)

const (
	Port = ":8080"
)

func main() {
	server := api.NewServer(Port)
	log.Printf("Starting server at port %s", Port)
	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", Port)
	}
}
