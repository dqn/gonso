package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
)

type MessageRequest struct {
	Body string `json:"body"`
	Type string `json:"type"`
}

func randomBytes(n uint) []byte {
	letter := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]byte, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return b
}

func postMessage(token string, messageRequest *MessageRequest) error {
	rawJSON, err := json.Marshal(messageRequest)
	if err != nil {
		return err
	}

	rawURL := "https://web.sd.lp1.acbaa.srv.nintendo.net/api/sd/v1/messages"
	req, err := http.NewRequest("POST", rawURL, bytes.NewBuffer([]byte(rawJSON)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}

func run() error {
	token := ""
	err := postMessage(token, &MessageRequest{
		Body: string(randomBytes(20)),
		Type: "all_friend",
	})

	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}
