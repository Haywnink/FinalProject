package api

import (
	"net/http"
	"strconv"
	"time"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	t, err := database.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}
	if t.Repeat == "" {
		if err := database.DeleteTask(id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not delete task"})
			return
		}
	} else {
		loc := time.Now().Location()
		baseDate, err := time.ParseInLocation("20060102", t.Date, loc)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "invalid task date"})
			return
		}
		next, err := NextDate(baseDate, t.Date, t.Repeat)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		t.Date = next
		if err := database.UpdateTask(t); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not update task"})
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{})
}
