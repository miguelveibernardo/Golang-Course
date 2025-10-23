package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"todo-cli/list"
)

func handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	items := list.LoadFromFile(list.DefaultDataFile)
	items = list.Add(items, body.Description)
	list.SaveToFile(list.DefaultDataFile, items)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	items := list.LoadFromFile(list.DefaultDataFile)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		ID    int    `json:"id"`
		Field string `json:"field"`
		Value string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	items := list.LoadFromFile(list.DefaultDataFile)
	var err error

	switch body.Field {
	case "description":
		items, err = list.UpdateDescription(items, body.ID, body.Value)
	case "status":
		items, err = list.UpdateStatus(items, body.ID, body.Value)
	default:
		http.Error(w, "Invalid field(must be 'description' or 'status')", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	list.SaveToFile(list.DefaultDataFile, items)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)

}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	items := list.LoadFromFile(list.DefaultDataFile)
	items = list.Delete(items, id)
	list.SaveToFile(list.DefaultDataFile, items)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/create", handleCreate)
	mux.HandleFunc("/get", handleGet)
	mux.HandleFunc("/update", handleUpdate)
	mux.HandleFunc("/delete", handleDelete)

	slog.Info("Starting HTTP server", "port", 8080)
	http.ListenAndServe(":8080", mux)
}
