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
			return "", fmt.Errorf("invalid day format")
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil || n < 1 || n > 400 {
			return "", fmt.Errorf("invalid day interval")
		}
		date := orig
		for {
			date = date.AddDate(0, 0, n)
			if date.After(now) {
				break
			}
		}
		return date.Format(dateLayout), nil

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("invalid year format")
		}
		date := orig
		for {
			date = date.AddDate(1, 0, 0)
			if orig.Month() == 2 && orig.Day() == 29 && date.Month() == 2 && date.Day() == 28 {
				date = time.Date(date.Year(), 3, 1, 0, 0, 0, 0, orig.Location())
			}
			if date.After(now) {
				break
			}
		}
		return date.Format(dateLayout), nil

	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid weekly format")
		}
		for _, ds := range strings.Split(parts[1], ",") {
			w, err := strconv.Atoi(ds)
			if err != nil || w < 1 || w > 7 {
				return "", fmt.Errorf("invalid weekly format")
			}
		}
		return nextWeekly(now, parts[1]).Format(dateLayout), nil

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("invalid monthly format")
		}
		var dayInts []int
		for _, ds := range strings.Split(parts[1], ",") {
			dv, err := strconv.Atoi(ds)
			if err != nil || dv == 0 || abs(dv) > 31 {
				return "", fmt.Errorf("invalid day interval")
			}
			dayInts = append(dayInts, dv)
		}
		var negAbs []int
		for _, v := range dayInts {
			if v < 0 {
				negAbs = append(negAbs, -v)
			}
		}
		if len(negAbs) > 1 {
			maxAbs := 0
			for _, a := range negAbs {
				if a > maxAbs {
					maxAbs = a
				}
			}
			got := map[int]bool{}
			for _, a := range negAbs {
				got[a] = true
			}
			for i := 1; i <= maxAbs; i++ {
				if !got[i] {
					return "", fmt.Errorf("invalid negative day offsets")
				}
			}
		}
		monthSet := map[int]bool{}
		if len(parts) == 3 {
			for _, ms := range strings.Split(parts[2], ",") {
				mi, err := strconv.Atoi(ms)
				if err != nil || mi < 1 || mi > 12 {
					return "", fmt.Errorf("invalid month token")
				}
				monthSet[mi] = true
			}
		}
		refDate := now
		if orig.After(now) {
			refDate = orig
		}
		date := refDate
		dimFn := func(y int, m time.Month) int {
			t := time.Date(y, m+1, 0, 0, 0, 0, 0, orig.Location())
			return t.Day()
		}
		for {
			date = date.AddDate(0, 0, 1)
			d := date.Day()
			m := int(date.Month())
			if len(monthSet) > 0 && !monthSet[m] {
				continue
			}
			ok := false
			dim := dimFn(date.Year(), date.Month())
			for _, dv := range dayInts {
				if dv > 0 && d == dv {
					ok = true
					break
				}
				if dv < 0 && d == dim+dv+1 {
					ok = true
					break
				}
			}
			if ok {
				return date.Format(dateLayout), nil
			}
		}

	default:
		return "", fmt.Errorf("unsupported repeat format")
	}
}

func nextWeekly(now time.Time, tok string) time.Time {
	allowed := map[int]bool{}
	for _, d := range strings.Split(tok, ",") {
		w, _ := strconv.Atoi(d)
		allowed[w] = true
	}
	date := now
	for {
		date = date.AddDate(0, 0, 1)
		wd := int(date.Weekday())
		if wd == 0 {
			wd = 7
		}
		if allowed[wd] {
			return date
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
