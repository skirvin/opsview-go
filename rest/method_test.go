package rest

import (
	"testing"
)

func TestMethod(t *testing.T) {
	var m Method
	m.String()

	if m.String() != "GET" {
		t.Errorf("Method is incorrect, got: %s, expected: %s.", m.String(), "GET")
	}
	if PATCH.String() != "PATCH" {
		t.Errorf("Method is incorrect, got: %s, expected: %s.", PATCH.String(), "PATCH")
	}
	if POST.String() != "POST" {
		t.Errorf("Method is incorrect, got: %s, expected: %s.", POST.String(), "POST")
	}
	if PUT.String() != "PUT" {
		t.Errorf("Method is incorrect, got: %s, expected: %s.", PUT.String(), "PUT")
	}
	if DELETE.String() != "DELETE" {
		t.Errorf("Method is incorrect, got: %s, expected: %s.", DELETE.String(), "DELETE")
	}
}
