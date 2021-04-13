package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler

	h := NoSurf(&myH)

	// Check return type
	switch v := h.(type) {
	case http.Handler:
		// do nothing, as expected
	default:
		t.Error("type is not http.Handler", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler

	h := SessionLoad(&myH)

	// Check return type
	switch v := h.(type) {
	case http.Handler:
		// do nothing, as expected
	default:
		t.Error("type is not http.Handler", v)
	}
}
