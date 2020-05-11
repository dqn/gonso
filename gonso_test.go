package gonso

import "testing"

func TestALL(t *testing.T) {
	n := New()
	_, err := n.Auth()
	if err != nil {
		t.Fatal(err)
	}
}
