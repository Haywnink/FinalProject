package main

import (
	"github.com/Haywnink/FinalProject/pkg/server"
	"log"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
