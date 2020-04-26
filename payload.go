package nso

type sessionTokenResponse struct {
	Code         string `json:"code"`
	SessionToken string `json:"session_token"`
	Error        string `json:"error"`
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
	Error       string   `json:"error"`
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
	Hash  string `json:"hash"`
	Error string `json:"error"`
}

type flagpResult struct {
	F  string `json:"f"`
	P1 string `json:"p1"`
	P2 string `json:"p2"`
	P3 string `json:"p3"`
}

type flagpResponse struct {
	Result flagpResult `json:"result"`
	Error  string      `json:"error"`
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
	Error         string              `json:"error"`
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
	Error         string                        `json:"error"`
}
