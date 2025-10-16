package api

import (
	"net/http"

	"github.com/magabrotheeeer/Yandex-Practicum-Final-Project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, "Задача не найдена", http.StatusNotFound)
		return
	}
	writeJSON(w, task)
}
