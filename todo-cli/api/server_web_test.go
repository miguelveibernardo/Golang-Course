package api_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"todo-cli/api"
)

func getTestMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/about", api.HandleAbout)
	mux.HandleFunc("/list", api.HandleListPage)
	return mux
}

// test the /about endpoint (static html)
func TestAboutPage(t *testing.T) {
	mux := getTestMux()
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	body, _ := io.ReadAll(w.Body)
	html := string(body)

	if !strings.Contains(html, "<h1>About ToDo Cli</h1>") {
		t.Errorf("Expected About page header missing, got: \n%s", html)
	}
}

// test the /list endpoint (dynamic html)
func TestListPage(t *testing.T) {
	mux := getTestMux()
	req := httptest.NewRequest(http.MethodGet, "/list", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	body, _ := io.ReadAll(w.Body)
	html := string(body)

	//verifying some basic HTML structure
	if !strings.Contains(html, "<h1>To-Do List</h1>") {
		t.Errorf("Expected List page header missing, got: \n%s", html)
	}

	if !strings.Contains(html, "No to-do items") && !strings.Contains(html, "<table>") {
		t.Errorf("Expected to see either a table or a empty message, got \n%s", html)
	}
}
