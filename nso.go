package nso

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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
	q.Set("state", state)
	q.Set("redirect_uri", "npf71b963c1b7b6d119://auth")
	q.Set("client_id", clientID)
	q.Set("scope", "openid user user.birthday user.mii user.screenName")
	q.Set("response_type", "session_token_code")
	q.Set("session_token_code_challenge", sessionTokenCodeChallenge)
	q.Set("session_token_code_challenge_method", "S256")
	q.Set("theme", "login_form")
	u.RawQuery = q.Encode()

	return u.String()
}

func processRequest(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	println(string(b))

	return b, nil
}

func postJSON(url string, header *http.Header, body interface{}) ([]byte, error) {
	var buf *bytes.Buffer
	if body != nil {
		rawJSON, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(rawJSON)
	}

	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}

	if header != nil {
		req.Header = *header
	}
	req.Header.Set("content-type", "application/json; charset=utf-8")

	b, err := processRequest(req)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func postForm(url string, header *http.Header, body *url.Values) ([]byte, error) {
	var reader io.Reader
	if body != nil {
		reader = strings.NewReader(body.Encode())
	}

	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}

	if header != nil {
		req.Header = *header
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	b, err := processRequest(req)
	if err != nil {
		return nil, err
	}

	return b, err
}

func getSessionToken(sessionTokenCode, sessionTokenCodeVerifier string) (*sessionTokenResponse, error) {
	u := "https://accounts.nintendo.com/connect/1.0.0/api/session_token"
	body := &url.Values{
		"client_id":                   {clientID},
		"session_token_code":          {sessionTokenCode},
		"session_token_code_verifier": {sessionTokenCodeVerifier},
	}
	b, err := postForm(u, nil, body)
	if err != nil {
		return nil, err
	}

	var r sessionTokenResponse
	err = json.Unmarshal(b, &r)

	return &r, err
}

func getToken(sessionToken string) (*tokenResponse, error) {
	u := "https://accounts.nintendo.com/connect/1.0.0/api/token"
	body := &tokenRequest{
		clientID,
		"urn:ietf:params:oauth:grant-type:jwt-bearer-session-token",
		sessionToken,
	}

	b, err := postJSON(u, nil, body)
	if err != nil {
		return nil, err
	}

	var r tokenResponse
	err = json.Unmarshal(b, &r)

	return &r, err
}

func callS2SAPI(naIDToken string, timestamp int64) (*s2sResponse, error) {
	u := "https://elifessler.com/s2s/api/gen2"
	header := &http.Header{
		"User-Agent": {"user_agent/version.num"},
	}
	body := &url.Values{
		"naIdToken": {naIDToken},
		"timestamp": {strconv.FormatInt(timestamp, 10)},
	}

	b, err := postForm(u, header, body)
	if err != nil {
		return nil, err
	}

	var r s2sResponse
	err = json.Unmarshal(b, &r)

	return &r, err
}

func callFlapgAPI(iid, idToken, guid, hash string, timestamp int64) (*flagpResponse, error) {
	u := "https://flapg.com/ika2/api/login?public"
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"x-token": {idToken},
		"x-time":  {strconv.FormatInt(timestamp, 10)},
		"x-guid":  {guid},
		"x-hash":  {hash},
		"x-ver":   {"3"},
		"x-iid":   {iid},
	}

	b, err := processRequest(req)
	if err != nil {
		return nil, err
	}

	var r flagpResponse
	err = json.Unmarshal(b, &r)

	return &r, err
}

func login(idToken, f, guid string, timestamp int64) (*loginResponse, error) {
	u := "https://api-lp1.znc.srv.nintendo.net/v1/Account/Login"
	header := &http.Header{
		"x-productversion": {"1.6.1.2"},
		"x-platform":       {"Android"},
	}
	body := &loginRequest{
		loginRequestParameter{
			f,
			"ja-JP",
			"1998-10-06",
			"JP",
			idToken,
			guid,
			timestamp,
		},
	}

	b, err := postJSON(u, header, body)
	if err != nil {
		return nil, err
	}

	var r loginResponse
	err = json.Unmarshal(b, &r)

	return &r, err
}

func getWebServiseToken(accessToken, f, registrationToken, guid string, timestamp int64) (*webServiceTokenResponse, error) {
	u := "https://api-lp1.znc.srv.nintendo.net/v2/Game/GetWebServiceToken"
	header := &http.Header{
		"authorization":    {fmt.Sprintf("Bearer %s", accessToken)},
		"x-productversion": {"1.6.1.2"},
		"x-platform":       {"Android"},
	}
	body := &webServiceTokenRequest{
		webServiceTokenRequestParameter{
			f,
			4953919198265344,
			registrationToken,
			guid,
			timestamp,
		},
	}

	b, err := postJSON(u, header, body)
	if err != nil {
		return nil, err
	}

	var r webServiceTokenResponse
	err = json.Unmarshal(b, &r)

	return &r, err
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
