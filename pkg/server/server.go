package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Haywnink/FinalProject/pkg/api"
	"github.com/Haywnink/FinalProject/pkg/db"
)

func Run() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	if err := db.Init("scheduler.db"); err != nil {
		return fmt.Errorf("ошибка инициализации БД: %v", err)
	}
	api.Init()
	http.Handle("/", http.FileServer(http.Dir("web")))
	fmt.Printf("Слушаем порт :%s\n", port)
	return http.ListenAndServe(":"+port, nil)
}
