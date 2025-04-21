package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("пустое правило")
	}
	loc := now.Location()
	orig, err := time.ParseInLocation("20060102", dstart, loc)
	if err != nil {
		return "", fmt.Errorf("некорректная дата")
	}

	parts := strings.Split(repeat, " ")
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("некорректный формат")
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil || n < 1 || n > 400 {
			return "", fmt.Errorf("некорректный интервал дней")
		}
		date := orig.AddDate(0, 0, n)
		for !date.After(now) {
			date = date.AddDate(0, 0, n)
		}
		return date.Format("20060102"), nil

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("некорректный формат")
		}
		date := orig.AddDate(1, 0, 0)
		if orig.Month() == 2 && orig.Day() == 29 &&
			date.Month() == 2 && date.Day() == 28 {
			date = time.Date(date.Year(), 3, 1, 0, 0, 0, 0, loc)
		}
		for !date.After(now) {
			date = date.AddDate(1, 0, 0)
			if orig.Month() == 2 && orig.Day() == 29 &&
				date.Month() == 2 && date.Day() == 28 {
				date = time.Date(date.Year(), 3, 1, 0, 0, 0, 0, loc)
			}
		}
		return date.Format("20060102"), nil

	default:
		return "", fmt.Errorf("неподдерживаемый формат")
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.ParseInLocation("20060102", nowStr, time.Local)
		if err != nil {
			http.Error(w, "некорректная дата now", http.StatusBadRequest)
			return
		}
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
