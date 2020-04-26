# go-nso

Nintendo Switch Online API wrapper

## Installation

```bash
$ go get github.com/dqn/go-nso
```

## Usage

```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dqn/go-nso"
)

type ACNHUsers struct {
	Users []struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Image string `json:"image"`
		Land  struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			DisplayID int    `json:"displayId"`
		} `json:"land"`
	} `json:"users"`
}

func main() {
	n := nso.New()
	accessToken, err := n.Auth()
	if err != nil {
		panic("failed to authenticate")
	}

	// Example for Animal Crossing: New Horizons API
	url := "https://web.sd.lp1.acbaa.srv.nintendo.net/api/sd/v1/users"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", fmt.Sprintf("_gtoken=%s", accessToken))

	resp, _ := http.DefaultClient.Do(req)
	b, _ := ioutil.ReadAll(resp.Body)

	var a ACNHUsers
	json.Unmarshal(b, &a)

	user := a.Users[0]
	fmt.Println(user.Name, user.Land.Name) // => どきゅん プリズム
}
```

First time, you need to authenticate.

```bash
authenticate by visiting this url: https://accounts.nintendo.com/connect/1.0.0/authorize?xxx=xxx
session token code: <input your session token code>
```

Credentials are cached as `./nso.json`.
