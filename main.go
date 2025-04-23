package main

import (
	"log"

	"github.com/Haywnink/FinalProject/pkg/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
