package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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

func parseCookie(s string) []*http.Cookie {
	// rawCookies := "cookie1=value1;cookie2=value2"

	h := http.Header{}
	h.Add("Cookie", s)
	r := http.Request{Header: h}
	return r.Cookies()
}

func run() error {
	rawURL := "https://web.sd.lp1.acbaa.srv.nintendo.net/api/sd/v1/messages"
	u, err := url.Parse("https://web.sd.lp1.acbaa.srv.nintendo.net/")
	if err != nil {
		return err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	rawCookie := ""
	jar.SetCookies(u, parseCookie(rawCookie))

	rawJSON, err := json.Marshal(&MessageRequest{
		Body: string(randomBytes(20)),
		Type: "all_friend",
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", rawURL, bytes.NewBuffer([]byte(rawJSON)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "")
	client := http.Client{Jar: jar}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(resp.Header)
	fmt.Println(resp.Header.Get("Set-Cookie"))

	fmt.Println(string(b))
	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}
