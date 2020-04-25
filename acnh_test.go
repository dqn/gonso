package acnh

import "testing"

func TestPostMessage(t *testing.T) {
	token := ""
	err := postMessage(token, &MessageRequest{
		Body: string(randomBytes(20)),
		Type: "all_friend",
	})
	if err != nil {
		t.Fatal(err)
	}
}
