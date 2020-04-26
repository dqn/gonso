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

	"github.com/google/uuid"
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

type loginRequestParameter struct {
	F          string `json:"f"`
	Language   string `json:"language"`
	NaBirthday string `json:"naBirthday"`
	NaCountry  string `json:"naCountry"`
	NaIDToken  string `json:"naIdToken"`
	RequestID  string `json:"requestId"`
	Timestamp  int64  `json:"timestamp"`
}

type loginRequest struct {
	Parameter loginRequestParameter `json:"parameter"`
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

type s2sResponse struct {
	Hash string `json:"hash"`
}

type flagpResult struct {
	F  string `json:"f"`
	P1 string `json:"p1"`
	P2 string `json:"p2"`
	P3 string `json:"p3"`
}

type flagpResponse struct {
	Result flagpResult `json:"result"`
}

type loginResponseResult struct {
	FirebaseCredential     firebaseCredential     `json:"firebaseCredential"`
	User                   user                   `json:"user"`
	WebAPIServerCredential webApiServerCredential `json:"webApiServerCredential"`
}

type loginResponse struct {
	CorrelationID string              `json:"correlationId"`
	Result        loginResponseResult `json:"result"`
	Status        int                 `json:"status"`
}

type webServiceTokenRequestParameter struct {
	F                 string `json:"f"`
	ID                int64  `json:"id"`
	RegistrationToken string `json:"registrationToken"`
	RequestID         string `json:"requestId"`
	Timestamp         int64  `json:"timestamp"`
}

type webServiceTokenRequest struct {
	Parameter webServiceTokenRequestParameter `json:"parameter"`
}

type webServiceTokenResponseResult struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
}

type webServiceTokenResponse struct {
	CorrelationID string                        `json:"correlationId"`
	Result        webServiceTokenResponseResult `json:"result"`
	Status        int                           `json:"status"`
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

func getSessionToken(sessionTokenCode, sessionTokenCodeVerifier string) (*sessionTokenResponse, error) {
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

func getToken(sessionToken string) (*tokenResponse, error) {
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

func callS2SAPI(naIDToken string, timestamp int64) (*s2sResponse, error) {
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
	println(string(b))

	return &s, err
}

func callFlapgAPI(iid, idToken, guid, hash string, timestamp int64) (*flagpResponse, error) {
	rawURL := "https://flapg.com/ika2/api/login?public"
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-token", idToken)
	req.Header.Set("x-time", strconv.FormatInt(timestamp, 10))
	req.Header.Set("x-guid", guid)
	req.Header.Set("x-hash", hash)
	req.Header.Set("x-ver", "3")
	req.Header.Set("x-iid", iid)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var t flagpResponse
	json.Unmarshal(b, &t)
	println(string(b))

	return &t, nil
}

func login(idToken, f, guid string, timestamp int64) (*loginResponse, error) {
	rawURL := "https://api-lp1.znc.srv.nintendo.net/v1/Account/Login"
	rawJSON, err := json.Marshal(loginRequest{
		loginRequestParameter{
			f,
			"ja-JP",
			"1998-10-06",
			"JP",
			idToken,
			guid,
			timestamp,
		},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", rawURL, bytes.NewBuffer(rawJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/json; charset=utf-8")
	req.Header.Set("x-productversion", "1.6.1.2")
	req.Header.Set("x-platform", "Android")

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
	json.Unmarshal(b, &l)
	println(string(b))

	return &l, nil
}

func getWebServiseToken(accessToken, f, registrationToken, guid string, timestamp int64) (*webServiceTokenResponse, error) {
	rawURL := "https://api-lp1.znc.srv.nintendo.net/v2/Game/GetWebServiceToken"
	rawJSON, err := json.Marshal(webServiceTokenRequest{
		webServiceTokenRequestParameter{
			f,
			4953919198265344,
			registrationToken,
			guid,
			timestamp,
		},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", rawURL, bytes.NewBuffer(rawJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/json; charset=utf-8")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("x-productversion", "1.6.1.2")
	req.Header.Set("x-platform", "Android")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var w webServiceTokenResponse
	json.Unmarshal(b, &w)
	println(string(b))

	return &w, nil
}

func (n *NSO) Auth() error {
	rand.Seed(time.Now().UnixNano())
	state := safeBase64Encode(randomBytes(36))
	sessionTokenCodeVerifier := safeBase64Encode(randomBytes(32))
	hash := sha256.Sum256([]byte(sessionTokenCodeVerifier))
	sessionTokenCodeChallenge := safeBase64Encode(hash[:])
	u := generateAuthURL(state, sessionTokenCodeChallenge)

	fmt.Printf("authorize by visiting this url: %s\n", u)

	var sessionTokenCode string
	fmt.Print("session token code: ")
	fmt.Scanf("%s", &sessionTokenCode)

	st, err := getSessionToken(sessionTokenCode, sessionTokenCodeVerifier)
	if err != nil {
		return err
	}

	t, err := getToken(st.SessionToken)
	if err != nil {
		return err
	}

	timestamp := time.Now().Unix()
	h, err := callS2SAPI(t.IDToken, timestamp)
	if err != nil {
		return err
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	guid := uuid.String()
	r, err := callFlapgAPI("nso", t.IDToken, guid, h.Hash, timestamp)
	if err != nil {
		return err
	}

	l, err := login(t.IDToken, r.Result.F, guid, timestamp)
	if err != nil {
		return err
	}

	accessToken := l.Result.WebAPIServerCredential.AccessToken
	h, err = callS2SAPI(accessToken, timestamp)
	if err != nil {
		return err
	}

	r, err = callFlapgAPI("app", accessToken, guid, h.Hash, timestamp)
	if err != nil {
		return err
	}

	w, err := getWebServiseToken(accessToken, r.Result.F, r.Result.P1, r.Result.P3, timestamp)
	if err != nil {
		return err
	}

	fmt.Println(w)

	return nil
}
