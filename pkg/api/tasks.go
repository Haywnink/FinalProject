package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func tasksHandler(w http.ResponseWriter, r *http.Request) {

	list, err := database.Tasks(50)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not get tasks"})
		return
	}

	search := r.URL.Query().Get("search")

	now := time.Now()
	loc := now.Location()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	out := make([]map[string]string, 0, len(list))

	for _, t := range list {
		rawDate := t.Date

		d, err := time.ParseInLocation(dateFormat, rawDate, loc)
		if err != nil {
			continue
		}

		if d.Before(today) {
			if t.Repeat != "" {
				next, err := NextDate(now, rawDate, t.Repeat)
				if err == nil {
					rawDate = next
				} else {
					rawDate = today.Format(dateFormat)
				}
			} else {
				rawDate = today.Format(dateFormat)
			}
		}

		if search != "" {
			formattedDate := d.Format("02.01.2006")
			if !strings.Contains(t.Title, search) &&
				!strings.Contains(t.Comment, search) &&
				!strings.Contains(formattedDate, search) {
				continue
			}
		}

		m := map[string]string{
			"id":      fmt.Sprint(t.ID),
			"date":    rawDate,
			"title":   t.Title,
			"comment": t.Comment,
			"repeat":  t.Repeat,
		}
		out = append(out, m)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"tasks": out})
}
