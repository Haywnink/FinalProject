package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Haywnink/FinalProject/pkg/api"
	"github.com/Haywnink/FinalProject/pkg/db"
)

func Run() error {
	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		dbPath = "scheduler.db"
	}
	database, err := db.New(dbPath)
	if err != nil {
		return fmt.Errorf("database initialization error: %v", err)
	}
	defer database.Close()

	api.Init(database)

	http.Handle("/", http.FileServer(http.Dir("web")))
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	fmt.Printf("Listening on port %s\n", port)
	return http.ListenAndServe(":"+port, nil)
}
