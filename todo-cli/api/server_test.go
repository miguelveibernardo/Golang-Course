package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-cli/api"
)

func getMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/create", api.HandleCreate)
	mux.HandleFunc("/get", api.HandleGet)
	mux.HandleFunc("/update", api.HandleUpdate)
	mux.HandleFunc("/delete", api.HandleDelete)
	return mux
}

func TestCreateItem(t *testing.T) {
	mux := getMux()
	req := httptest.NewRequest(http.MethodPost, "/create?description=TestTask", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	var items []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &items); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if len(items) == 0 {
		t.Fatalf("Expected at least one item in response, got %d", len(items))
	}

}

func TestGetItems(t *testing.T) {
	mux := getMux()
	req := httptest.NewRequest(http.MethodGet, "/get", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	var items []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &items); err != nil {
		t.Fatalf("Fail to parse responde: %v", err)
	}
}

func TestUpdateInvalidItem(t *testing.T) {
	mux := getMux()
	req := httptest.NewRequest(http.MethodPut, "/update?id=999&field=description&value=Nope", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 Bad Request, got %d", w.Code)
	}
}

func TestDeleteInvalidItem(t *testing.T) {
	mux := getMux()
	req := httptest.NewRequest(http.MethodDelete, "/delete?id=999", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}
}
