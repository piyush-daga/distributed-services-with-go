package main

import (
	"log"

	"github.com/piyush-daga/proglog/internal/server"
)

func main() {
	s := server.NewHTTPServer(":8080")
	log.Fatal(s.ListenAndServe())
}
