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

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.ReqTask
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJSONError(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if t.ID == "" {
		writeJSONError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}
	if err := validateTask(&t); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.UpdateTask(&t); err != nil {
		writeJSONError(w, "Задача не найдена", http.StatusNotFound)
		return
	}
	writeJSON(w, map[string]string{})
}

func validateTask(t *db.ReqTask) error {
	dt, err := time.Parse("20060102", t.Date)
	if err != nil {
		return fmt.Errorf("invalid date")
	}
	today := time.Now().Truncate(24 * time.Hour)
	if dt.Before(today) {
		return fmt.Errorf("invalid date")
	}
	if strings.TrimSpace(t.Title) == "" {
		return fmt.Errorf("invalid title")
	}
	if t.Repeat != "" {
		parts := strings.SplitN(t.Repeat, " ", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid repeat")
		}
		if parts[0] != "d" && parts[0] != "w" {
			return fmt.Errorf("invalid repeat")
		}
		if _, err := strconv.Atoi(parts[1]); err != nil {
			return fmt.Errorf("invalid repeat")
		}
	}
	return nil
}
