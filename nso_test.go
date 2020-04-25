package nso

import "testing"

func TestLogin(t *testing.T) {
	n := New()
	err := n.Login("", "")
	if err != nil {
		t.Fatal(err)
	}
}
