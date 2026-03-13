package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHome(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	req := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	body := string(bodyBytes)
	if !strings.Contains(body, "/health") || !strings.Contains(body, "/api/sessions") {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestHome_NotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	req := httptest.NewRequest(http.MethodGet, "http://example.com/not-exist", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, res.StatusCode)
	}
}
