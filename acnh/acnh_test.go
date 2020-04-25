package acnh

import (
	"math/rand"
	"testing"
)

func randomBytes(n uint) []byte {
	letter := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]byte, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return b
}
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
