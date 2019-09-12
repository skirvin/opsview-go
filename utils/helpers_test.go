package utils

import (
	"testing"
)

func TestSanatize(t *testing.T) {
	if Sanatize("http://www.example/com/test/") != "http://www.example.com/test" {
		t.Errorf("Sanatize result is incorrect, got: %s, expected: %s.", "http://www.example/com/test/", "http://www.example/com/test")
	}
}
