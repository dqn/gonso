package acnh

import (
	"testing"
)

func TestPostMessage(t *testing.T) {
	token := ""
	err := postMessage(token, &MessageRequest{
		Body: "Test Message",
		Type: "all_friend",
	})
	if err != nil {
		t.Fatal(err)
	}
}
