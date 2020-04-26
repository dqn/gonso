package nso

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

type sessionTokenResponse struct {
	Code         string `json:"code"`
	SessionToken string `json:"session_token"`
}

type tokenRequest struct {
	ClientID     string `json:"client_id"`
	GrantType    string `json:"grant_type"`
	SessionToken string `json:"session_token"`
}

type tokenResponse struct {
	AccessToken string   `json:"access_token"`
	ExpiresIn   uint     `json:"expires_in"`
	IDToken     string   `json:"id_token"`
	Scope       []string `json:"scope"`
	TokenType   string   `json:"token_type"`
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

func fetchSessionToken(sessionTokenCode, sessionTokenCodeVerifier string) (*sessionTokenResponse, error) {
	rawURL := "https://accounts.nintendo.com/connect/1.0.0/api/session_token"
	values := &url.Values{}
	values.Set("client_id", clientID)
	values.Set("session_token_code", sessionTokenCode)
	values.Set("session_token_code_verifier", sessionTokenCodeVerifier)

	req, err := http.NewRequest("POST", rawURL, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var st sessionTokenResponse
	json.Unmarshal(b, &st)

	return &st, nil
}

func fetchToken(sessionToken string) (*tokenResponse, error) {
	rawURL := "https://accounts.nintendo.com/connect/1.0.0/api/token"
	rawJSON, err := json.Marshal(tokenRequest{
		clientID,
		"urn:ietf:params:oauth:grant-type:jwt-bearer-session-token",
		sessionToken,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", rawURL, bytes.NewBuffer(rawJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var t tokenResponse
	json.Unmarshal(b, &t)

	return &t, nil
}

func (n *NSO) Auth() error {
	state := safeBase64Encode(randomBytes(36))
	sessionTokenCodeVerifier := safeBase64Encode(randomBytes(32))
	hash := sha256.Sum256([]byte(sessionTokenCodeVerifier))
	sessionTokenCodeChallenge := safeBase64Encode(hash[:])
	u := generateAuthURL(state, sessionTokenCodeChallenge)

	fmt.Printf("authorize by visiting this url: %s\n", u)

	var sessionTokenCode string
	fmt.Print("session token code: ")
	fmt.Scanf("%s", &sessionTokenCode)

	st, err := fetchSessionToken(sessionTokenCode, sessionTokenCodeVerifier)
	if err != nil {
		return err
	}
	t, err := fetchToken(st.SessionToken)
	if err != nil {
		return err
	}
	fmt.Println(t)

	return nil
}
