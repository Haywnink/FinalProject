package api

import (
	"net/http"

	"github.com/Haywnink/FinalProject/pkg/db"
)

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	list, err := db.Tasks(50)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка получения списка задач"})
		return
	}
	writeJSON(w, map[string]interface{}{"tasks": list})
}
