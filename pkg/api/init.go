package api

import "net/http"

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDateHandler)
	mux.HandleFunc("/api/task", processRequest)
	mux.HandleFunc("/api/tasks", readTasks)
	mux.HandleFunc("/api/task/done", doneTaskHandler)
}

func processRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
