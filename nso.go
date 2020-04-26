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

const clientID = "71b963c1b7b6d119"
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type NSO struct {
	client *http.Client
}

func New() *NSO {
	return &NSO{&http.Client{}}
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return b
}

func generateAuthURL(state, sessionTokenCodeChallenge string) string {
	u, _ := url.Parse("https://accounts.nintendo.com/connect/1.0.0/authorize")
	q := &url.Values{
		"state":                               {state},
		"redirect_uri":                        {"npf71b963c1b7b6d119://auth"},
		"client_id":                           {clientID},
		"scope":                               {"openid user user.birthday user.mii user.screenName"},
		"response_type":                       {"session_token_code"},
		"session_token_code_challenge":        {sessionTokenCodeChallenge},
		"session_token_code_challenge_method": {"S256"},
		"theme":                               {"login_form"},
	}
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

	return processRequest(req)
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

	return processRequest(req)
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

func callS2SAPI(token string, timestamp int64) (*s2sResponse, error) {
	u := "https://elifessler.com/s2s/api/gen2"
	header := &http.Header{
		"User-Agent": {"user_agent/version.num"},
	}
	body := &url.Values{
		"naIdToken": {token},
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

// https://github.com/frozenpandaman/splatnet2statink/wiki/api-docs
func callFlapgAPI(iid, token, guid string, timestamp int64) (*flagpResponse, error) {
	h, err := callS2SAPI(token, timestamp)
	if err != nil {
		return nil, err
	}

	u := "https://flapg.com/ika2/api/login?public"
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"x-token": {token},
		"x-time":  {strconv.FormatInt(timestamp, 10)},
		"x-guid":  {guid},
		"x-hash":  {h.Hash},
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
	state := base64.RawURLEncoding.EncodeToString(randomBytes(36))
	sessionTokenCodeVerifier := base64.RawURLEncoding.EncodeToString(randomBytes(32))
	hash := sha256.Sum256([]byte(sessionTokenCodeVerifier))
	sessionTokenCodeChallenge := base64.RawURLEncoding.EncodeToString(hash[:])
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

	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	guid := uuid.String()
	timestamp := time.Now().Unix()

	r, err := callFlapgAPI("nso", t.IDToken, guid, timestamp)
	if err != nil {
		return err
	}

	l, err := login(t.IDToken, r.Result.F, guid, timestamp)
	if err != nil {
		return err
	}

	token := l.Result.WebAPIServerCredential.AccessToken
	r, err = callFlapgAPI("app", token, guid, timestamp)
	if err != nil {
		return err
	}

	w, err := getWebServiseToken(token, r.Result.F, r.Result.P1, guid, timestamp)
	if err != nil {
		return err
	}

	fmt.Println(w)

	return nil
}
