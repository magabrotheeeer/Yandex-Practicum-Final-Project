package api

import (
	"net/http"

	"github.com/magabrotheeeer/Yandex-Practicum-Final-Project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.ReqTask `json:"tasks"`
}

func readTasks(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tasks == nil {
		tasks = make([]*db.ReqTask, 0)
	}
	writeJSON(w, struct {
		Tasks []*db.ReqTask `json:"tasks"`
	}{Tasks: tasks})
}
