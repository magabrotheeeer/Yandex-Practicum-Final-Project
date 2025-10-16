package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/magabrotheeeer/Yandex-Practicum-Final-Project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSONError(w, "invalid json body", http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		writeJSONError(w, "title is required", http.StatusBadRequest)
		return
	}
	err = checkDate(&task)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{
		"id": strconv.FormatInt(id, 10),
	})
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("failed to parse field Date: %v", err)
	}

	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if strings.TrimSpace(task.Repeat) != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("repeat rule invalid: %v", err)
		}

		if t.Before(nowDate) {
			task.Date = next
		}
	} else {
		if t.Before(nowDate) {
			task.Date = now.Format("20060102")
		}
	}

	return nil
}
