package gonso

import (
	"fmt"
	"os"
	"testing"
)

func TestAuth(t *testing.T) {
	token, err := Auth(os.Getenv("SESSION_TOKEN"))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(token)
}
