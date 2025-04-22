package api

import (
	"net/http"

	"github.com/Haywnink/FinalProject/pkg/db"
)

var database *db.DB

func Init(dbInstance *db.DB) {
	database = dbInstance
	http.HandleFunc("/api/signin", SignInHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneHandler))
	http.HandleFunc("/api/nextdate", auth(nextDateHandler))
}
