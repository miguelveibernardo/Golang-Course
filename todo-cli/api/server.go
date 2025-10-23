package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"todo-cli/list"
)

func HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	description := r.URL.Query().Get("description")
	if description == "" {
		http.Error(w, "Missing description parameter", http.StatusBadRequest)
		return
	}

	items := list.LoadFromFile(list.DefaultDataFile)
	items = list.Add(items, description)
	list.SaveToFile(list.DefaultDataFile, items)

	slog.Info("Item created via API", "description", description)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	items := list.LoadFromFile(list.DefaultDataFile)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	field := r.URL.Query().Get("field")
	value := r.URL.Query().Get("value")

	if idStr == "" || field == "" || value == "" {
		http.Error(w, "Missing required parameters (id, field, value)", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	items := list.LoadFromFile(list.DefaultDataFile)

	switch field {
	case "description":
		items, err = list.UpdateDescription(items, id, value)
	case "status":
		items, err = list.UpdateStatus(items, id, value)
	default:
		http.Error(w, "Invalid field(must be 'description' or 'status')", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	list.SaveToFile(list.DefaultDataFile, items)
	slog.Info("Item updated via API", "id", id, "field", field, "value", value)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)

}

func HandleDelete(w http.ResponseWriter, r *http.Request) {
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

	slog.Info("Item deleted via API", "id", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/create", HandleCreate)
	mux.HandleFunc("/get", HandleGet)
	mux.HandleFunc("/update", HandleUpdate)
	mux.HandleFunc("/delete", HandleDelete)

	slog.Info("Starting HTTP server", "port", 8080)
	http.ListenAndServe(":8080", mux)
}
