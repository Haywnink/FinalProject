package api

import (
	"net/http"
)

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	list, err := database.Tasks(50)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not get tasks"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"tasks": list})
}
