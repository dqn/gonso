package nso

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var clientID = "71b963c1b7b6d119"

type NSO struct {
	client *http.Client
}

func New() *NSO {
	return &NSO{&http.Client{}}
}

func (n *NSO) Login(mail, password string) error {
	u := "https://accounts.nintendo.com/connect/1.0.0/api/session_token"
	sessionTokenCode := ""
	sessionTokenCodeVerifier := ""

	values := &url.Values{}
	values.Set("client_id", clientID)
	values.Set("session_token_code", sessionTokenCode)
	values.Set("session_token_code_verifier", sessionTokenCodeVerifier)

	req, err := http.NewRequest("POST", u, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	println(string(b))

	return nil
}
