package api

import "net/http"

func Init() {
	http.HandleFunc("/api/signin", SignInHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneHandler))
	http.HandleFunc("/api/nextdate", auth(nextDateHandler))
}
