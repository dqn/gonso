package acnh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type MessageRequest struct {
	Body string `json:"body"`
	Type string `json:"type"`
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
