package gonso

import (
	"os"
	"testing"
)

func TestAuth(t *testing.T) {
	_, err := Auth(os.Getenv("SESSION_TOKEN"))
	if err != nil {
		t.Fatal(err)
	}
}
