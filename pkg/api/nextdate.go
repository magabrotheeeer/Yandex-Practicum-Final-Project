package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/magabrotheeeer/Yandex-Practicum-Final-Project/pkg/db"
)

var (
	errEmptyRepeat   = errors.New("repeat: empty")
	errBadStartDate  = errors.New("dstart: invalid date")
	errBadRepeatForm = errors.New("repeat: invalid repeat")
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()
	dstart := q.Get("date")
	repeat := q.Get("repeat")
	if dstart == "" || repeat == "" {
		http.Error(w, "missing required query params: date, repeat", http.StatusBadRequest)
		return
	}

	var now time.Time
	if nowStr := q.Get("now"); nowStr == "" {
		now = time.Now()
	} else {
		t, err := time.ParseInLocation(db.DateFormat, nowStr, time.Local)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			http.Error(w, "now: invalid date", http.StatusBadRequest)
			return
		}
		now = t
	}

	next, err := NextDate(now, dstart, repeat)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(next))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", errEmptyRepeat
	}

	loc := now.Location()
	start, err := time.ParseInLocation(db.DateFormat, dstart, loc)
	if err != nil {
		return "", errBadStartDate
	}
	now = normalizeDate(now.In(loc))
	date := normalizeDate(start)

	parts := strings.Fields(repeat)
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errBadRepeatForm
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil || n < 1 || n > 400 {
			return "", errBadRepeatForm
		}
		date = date.AddDate(0, 0, n)
		for !afterNow(date, now) {
			date = date.AddDate(0, 0, n)
		}
		return date.Format(db.DateFormat), nil

	case "y":
		if len(parts) != 1 {
			return "", errBadRepeatForm
		}
		date = addYearHuman(date)
		for !afterNow(date, now) {
			date = addYearHuman(date)
		}
		return date.Format(db.DateFormat), nil

	case "w":
		if len(parts) != 2 {
			return "", errBadRepeatForm
		}
		wset, err := parseWeekdays(parts[1])
		if err != nil {
			return "", errBadRepeatForm
		}
		date = date.AddDate(0, 0, 1)
		for {
			if afterNow(date, now) && wset[weekdayISO(date)] {
				return date.Format(db.DateFormat), nil
			}
			date = date.AddDate(0, 0, 1)
			if date.Sub(start) > 100*365*24*time.Hour {
				return "", errors.New("repeat: search exceeded bounds")
			}
		}

	case "m":
		if len(parts) != 2 && len(parts) != 3 {
			return "", errBadRepeatForm
		}
		daySpec, err := parseMonthDays(parts[1])
		if err != nil {
			return "", errBadRepeatForm
		}
		monSpec := [13]bool{}
		if len(parts) == 3 {
			monSpec, err = parseMonths(parts[2])
			if err != nil {
				return "", errBadRepeatForm
			}
		} else {
			for i := 1; i <= 12; i++ {
				monSpec[i] = true
			}
		}
		date = date.AddDate(0, 0, 1)
		for {
			if afterNow(date, now) && monSpec[int(date.Month())] && matchMonthDay(date, daySpec) {
				return date.Format(db.DateFormat), nil
			}
			date = date.AddDate(0, 0, 1)
			if date.Sub(start) > 100*365*24*time.Hour {
				return "", errors.New("repeat: search exceeded bounds")
			}
		}

	default:
		return "", errBadRepeatForm
	}
}

func addYearHuman(prev time.Time) time.Time {
	t := prev.AddDate(1, 0, 0)
	if prev.Month() == time.February && prev.Day() == 29 &&
		t.Month() == time.February && t.Day() == 28 {
		t = t.AddDate(0, 0, 1)
	}
	return t
}

func normalizeDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func parseWeekdays(s string) ([8]bool, error) {
	var set [8]bool
	for _, tok := range splitCSV(s) {
		n, err := strconv.Atoi(tok)
		if err != nil || n < 1 || n > 7 {
			return set, errBadRepeatForm
		}
		set[n] = true
	}
	if !(set[1] || set[2] || set[3] || set[4] || set[5] || set[6] || set[7]) {
		return set, errBadRepeatForm
	}
	return set, nil
}

func weekdayISO(t time.Time) int {
	w := int(t.Weekday())
	if w == 0 {
		return 7
	}
	return w
}

type monthDaySpec struct {
	exact  [32]bool
	last   bool
	before bool
}

func parseMonthDays(s string) (monthDaySpec, error) {
	var spec monthDaySpec
	for _, tok := range splitCSV(s) {
		n, err := strconv.Atoi(tok)
		if err != nil {
			return spec, errBadRepeatForm
		}
		switch {
		case n >= 1 && n <= 31:
			spec.exact[n] = true
		case n == -1:
			spec.last = true
		case n == -2:
			spec.before = true
		default:
			return spec, errBadRepeatForm
		}
	}
	if emptyMonthDaySpec(spec) {
		return spec, errBadRepeatForm
	}
	return spec, nil
}

func emptyMonthDaySpec(s monthDaySpec) bool {
	if s.last || s.before {
		return false
	}
	for i := 1; i <= 31; i++ {
		if s.exact[i] {
			return false
		}
	}
	return true
}

func parseMonths(s string) ([13]bool, error) {
	var set [13]bool
	for _, tok := range splitCSV(s) {
		n, err := strconv.Atoi(tok)
		if err != nil || n < 1 || n > 12 {
			return set, errBadRepeatForm
		}
		set[n] = true
	}
	if !(set[1] || set[2] || set[3] || set[4] || set[5] || set[6] || set[7] || set[8] || set[9] || set[10] || set[11] || set[12]) {
		return set, errBadRepeatForm
	}
	return set, nil
}

func matchMonthDay(t time.Time, spec monthDaySpec) bool {
	d := t.Day()
	if d <= 31 && spec.exact[d] {
		return true
	}
	year, month, _ := t.Date()
	last := lastDayOfMonth(year, month)
	if spec.last && d == last {
		return true
	}
	if spec.before && d == last-1 {
		return true
	}
	return false
}

func lastDayOfMonth(y int, m time.Month) int {
	n := time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC)
	return n.Day()
}

func splitCSV(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return []string{}
	}
	raw := strings.Split(s, ",")
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		v = strings.TrimSpace(v)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}
