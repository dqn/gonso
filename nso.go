package nso

import (
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"math/rand"
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

func randomBytes(n uint) []byte {
	letter := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]byte, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return b
}

func safeBase64Encode(b []byte) string {
	s := base64.StdEncoding.EncodeToString(b)
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "+", "-")
	s = strings.ReplaceAll(s, "=", "")
	return s
}

func generateAuthURL(state, sessionTokenCodeChallenge string) string {
	u, _ := url.Parse("https://accounts.nintendo.com/connect/1.0.0/authorize")
	q := u.Query()
	q.Add("state", state)
	q.Add("redirect_uri", "npf71b963c1b7b6d119://auth")
	q.Add("client_id", clientID)
	q.Add("scope", "openid user user.birthday user.mii user.screenName")
	q.Add("response_type", "session_token_code")
	q.Add("session_token_code_challenge", sessionTokenCodeChallenge)
	q.Add("session_token_code_challenge_method", "S256")
	q.Add("theme", "login_form")
	u.RawQuery = q.Encode()

	return u.String()
}

func (n *NSO) Auth() string {
	state := safeBase64Encode(randomBytes(36))
	sessionTokenCodeVerifier := safeBase64Encode(randomBytes(32))
	hash := sha256.Sum256([]byte(sessionTokenCodeVerifier))
	sessionTokenCodeChallenge := safeBase64Encode(hash[:])
	u := generateAuthURL(state, sessionTokenCodeChallenge)

	println(u)

	return ""
}

func (n *NSO) fetchSessionToken(sessionTokenCode, sessionTokenCodeVerifier string) error {
	u := "https://accounts.nintendo.com/connect/1.0.0/api/session_token"

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
