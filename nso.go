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
	"strconv"
	"strings"
	"time"
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

type parameter struct {
	F          string `json:"f"`
	Language   string `json:"language"`
	NaBirthday string `json:"naBirthday"`
	NaCountry  string `json:"naCountry"`
	NaIDToken  string `json:"naIdToken"`
	RequestID  string `json:"requestId"`
	Timestamp  int64  `json:"timestamp"`
}

type loginRequest struct {
	Parameter parameter `json:"parameter"`
}

type firebaseCredential struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
}

type membership struct {
	Active bool `json:"active"`
}

type user struct {
	ID         int64      `json:"id"`
	ImageURI   string     `json:"imageUri"`
	Membership membership `json:"membership"`
	Name       string     `json:"name"`
	SupportID  string     `json:"supportId"`
}

type webApiServerCredential struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
}

type result struct {
	FirebaseCredential     firebaseCredential     `json:"firebaseCredential"`
	User                   user                   `json:"user"`
	WebAPIServerCredential webApiServerCredential `json:"webApiServerCredential"`
}

type s2sResponse struct {
	Hash string `json:"hash"`
}

type flagpResponse struct {
}

type loginResponse struct {
	CorrelationID string `json:"correlationId"`
	Result        result `json:"result"`
	Status        int    `json:"status"`
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

func fetchHashByS2S(naIDToken string, timestamp int64) (*s2sResponse, error) {
	rawURL := "https://elifessler.com/s2s/api/gen2"
	values := &url.Values{}
	values.Set("naIdToken", naIDToken)
	values.Set("timestamp", strconv.FormatInt(timestamp, 10))

	req, err := http.NewRequest("POST", rawURL, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "user_agent/version.num")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var s s2sResponse
	json.Unmarshal(b, &s)

	return &s, err
}

// func fetchRequestIdAndRequestID(idToken, hash string, unix int64) (*flagpResponse, error) {
// 	rawURL := "https://flapg.com/ika2/api/login?public"
// 	req, err := http.NewRequest("POST", rawURL, bytes.NewBuffer(rawJSON))
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("x-token", idToken)
// 	req.Header.Set("x-time", strconv.FormatInt(unix, 10))
// 	req.Header.Set("x-guid", "")
// 	req.Header.Set("x-hash", hash)
// 	req.Header.Set("x-ver", "3")
// 	req.Header.Set("x-iid", "nso")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	b, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var t tokenResponse
// 	json.Unmarshal(b, &t)

// 	return &t, nil
// }

func login(idToken string) (*loginResponse, error) {
	_, err := fetchHashByS2S(idToken, time.Now().Unix())
	if err != nil {
		return nil, err
	}
	return nil, nil
	rawURL := "https://api-lp1.znc.srv.nintendo.net/v1/Account/Login"
	rawJSON, err := json.Marshal(loginRequest{
		parameter{
			"1e3de1eedef9952d1eb7ecb6ae520fabc5828d9c9d2fef9f89c889e68b15358b40e109758f954d2c746511cf",
			"ja-JP",
			"1998-10-06",
			"JP",
			idToken,
			"28536744-db82-47e4-a30e-1d2dacbd1e24",
			time.Now().Unix(),
		},
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

	var l loginResponse
	println(string(b))
	json.Unmarshal(b, &l)

	return &l, nil
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
	l, err := login(t.IDToken)
	if err != nil {
		return err
	}
	fmt.Println(l)

	return nil
}
