package api

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"todo-cli/list"
)

// ensures that templates load correctly in both runtime and test modes(I was having trouble with test mode)
func resolvePath(relPath string) string {
	//Try relative to the current working directory
	if _, err := os.Stat(relPath); err == nil {
		return relPath
	}
	// then try one level up
	alt := filepath.Join("..", relPath)
	if _, err := os.Stat(alt); err == nil {
		return alt
	}
	return relPath
}

func HandleAbout(w http.ResponseWriter, r *http.Request) {
	path := resolvePath("web/about.html")
	http.ServeFile(w, r, path)
}

func HandleListPage(w http.ResponseWriter, r *http.Request) {
	items := list.LoadFromFile(list.DefaultDataFile)

	tmplpath := resolvePath("web/list.html")
	tmpl, err := template.ParseFiles(tmplpath)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		slog.Error("Template parse error", "error", err)
		return
	}

	data := struct {
		Items []list.Item
		Count int
	}{
		Items: items,
		Count: len(items),
	}

	w.Header().Set("Content-Type", "text/html; charset=u-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		slog.Error("Template execution error", "error", err)
	}
}

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

	traceID := GetTraceID(r.Context())
	slog.Info("Handling /post request", "trace_id", traceID)
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	items := list.LoadFromFile(list.DefaultDataFile)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)

	traceID := GetTraceID(r.Context())
	slog.Info("Handling /get request", "trace_id", traceID)
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

	traceID := GetTraceID(r.Context())
	slog.Info("Handling /put request", "trace_id", traceID)

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

	traceID := GetTraceID(r.Context())
	slog.Info("Handling /delete request", "trace_id", traceID)
}

func StartServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/create", HandleCreate)
	mux.HandleFunc("/get", HandleGet)
	mux.HandleFunc("/update", HandleUpdate)
	mux.HandleFunc("/delete", HandleDelete)

	//web routes
	mux.HandleFunc("/about", HandleAbout)
	mux.HandleFunc("/list", HandleListPage)

	handler := TraceMiddleware(mux)

	slog.Info("Starting HTTP server", "port", 8080)
	http.ListenAndServe(":8080", handler)
}
