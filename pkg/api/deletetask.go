package api

import (
	"net/http"

	"github.com/magabrotheeeer/Yandex-Practicum-Final-Project/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{})
}
