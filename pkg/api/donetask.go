package api

import (
	"net/http"
	"time"

	"github.com/magabrotheeeer/Yandex-Practicum-Final-Project/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	t, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, "Задача не найдена", http.StatusNotFound)
		return
	}
	if t.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		next, err := NextDate(time.Now(), t.Date, t.Repeat)
		if err != nil {
			writeJSONError(w, "не удалось вычислить следующую дату", http.StatusInternalServerError)
			return
		}
		if err := db.UpdateDate(next, id); err != nil {
			writeJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	writeJSON(w, map[string]string{})
}
