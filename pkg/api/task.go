package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Haywnink/FinalProject/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJSON(w, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}
	if t.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок"})
		return
	}

	now := time.Now()
	loc := now.Location()
	// начало сегодняшнего дня в локальной зоне
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	if t.Date == "" {
		// если нет даты — ставим сегодня
		t.Date = today.Format("20060102")
	} else {
		// парсим именно в локальной зоне
		d, err := time.ParseInLocation("20060102", t.Date, loc)
		if err != nil {
			writeJSON(w, map[string]string{"error": "некорректная дата"})
			return
		}
		if d.Before(today) {
			// если дата до начала сегодняшнего дня
			if t.Repeat != "" {
				// для повторяющихся задач — рассчитываем следующий запуск
				next, err := NextDate(now, t.Date, t.Repeat)
				if err != nil {
					writeJSON(w, map[string]string{"error": err.Error()})
					return
				}
				t.Date = next
			} else {
				// одноразовые задачи — просто ставим сегодня
				t.Date = today.Format("20060102")
			}
		}
	}

	id, err := db.AddTask(&t)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка добавления задачи"})
		return
	}
	writeJSON(w, map[string]string{"id": fmt.Sprint(id)})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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
	writeJSON(w, t)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJSON(w, map[string]string{"error": "ошибка разбора JSON"})
		return
	}
	if t.ID == "" {
		writeJSON(w, map[string]string{"error": "не указан идентификатор"})
		return
	}
	if t.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок"})
		return
	}

	now := time.Now()
	loc := now.Location()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	if t.Date == "" {
		t.Date = today.Format("20060102")
	} else {
		d, err := time.ParseInLocation("20060102", t.Date, loc)
		if err != nil {
			writeJSON(w, map[string]string{"error": "некорректная дата"})
			return
		}
		if d.Before(today) {
			if t.Repeat != "" {
				next, err := NextDate(now, t.Date, t.Repeat)
				if err != nil {
					writeJSON(w, map[string]string{"error": err.Error()})
					return
				}
				t.Date = next
			} else {
				t.Date = today.Format("20060102")
			}
		}
	}

	if err := db.UpdateTask(&t); err != nil {
		writeJSON(w, map[string]string{"error": "задача не найдена"})
		return
	}
	writeJSON(w, map[string]string{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "не указан идентификатор"})
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, map[string]string{"error": "ошибка удаления задачи"})
		return
	}
	writeJSON(w, map[string]string{})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}
