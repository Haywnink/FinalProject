package api

import (
	"net/http"
	"time"

	"github.com/Haywnink/FinalProject/pkg/db"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "не указан идентификатор"})
		return
	}
	t, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "задача не найдена"})
		return
	}
	if t.Repeat == "" {
		// одноразовая — удаляем
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]string{"error": "ошибка удаления задачи"})
			return
		}
	} else {
		// повторяющаяся — прибавляем ровно один интервал к хранимой дате
		loc := time.Now().Location()
		baseDate, err := time.ParseInLocation("20060102", t.Date, loc)
		if err != nil {
			writeJSON(w, map[string]string{"error": "некорректная дата задачи"})
			return
		}
		next, err := NextDate(baseDate, t.Date, t.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		if err := db.UpdateDate(next, id); err != nil {
			writeJSON(w, map[string]string{"error": "ошибка обновления даты"})
			return
		}
	}
	writeJSON(w, map[string]string{})
}
