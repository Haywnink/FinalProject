package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "20060102"

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("empty repeat rule")
	}
	orig, err := time.ParseInLocation(dateLayout, dstart, now.Location())
	if err != nil {
		return "", fmt.Errorf("invalid start date")
	}

	parts := strings.Split(repeat, " ")
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid format")
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil || n < 1 || n > 400 {
			return "", fmt.Errorf("invalid day interval")
		}
		date := orig.AddDate(0, 0, n)
		for !date.After(now) {
			date = date.AddDate(0, 0, n)
		}
		return date.Format(dateLayout), nil

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("invalid format")
		}
		date := orig.AddDate(1, 0, 0)
		if orig.Month() == 2 && orig.Day() == 29 &&
			date.Month() == 2 && date.Day() == 28 {
			date = time.Date(date.Year(), 3, 1, 0, 0, 0, 0, orig.Location())
		}
		for !date.After(now) {
			date = date.AddDate(1, 0, 0)
			if orig.Month() == 2 && orig.Day() == 29 &&
				date.Month() == 2 && date.Day() == 28 {
				date = time.Date(date.Year(), 3, 1, 0, 0, 0, 0, orig.Location())
			}
		}
		return date.Format(dateLayout), nil

	default:
		return "", fmt.Errorf("unsupported repeat format")
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	now := time.Now()
	if nowStr := r.URL.Query().Get("now"); nowStr != "" {
		parsed, err := time.ParseInLocation(dateLayout, nowStr, time.Local)
		if err != nil {
			http.Error(w, "invalid now date", http.StatusBadRequest)
			return
		}
		now = parsed
	}
	dstart := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")
	next, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(next))
}
