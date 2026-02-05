package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirectHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(redirectHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusTemporaryRedirect)
	}

	expected := "https://afonso.dev"
	if location := rr.Header().Get("Location"); location != expected {
		t.Errorf("handler returned wrong redirect location: got %v want %v",
			location, expected)
	}
}
