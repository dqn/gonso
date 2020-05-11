# gonso

Nintendo Switch Online API wrapper

## Installation

```bash
$ go get github.com/dqn/gonso
```

## Usage

```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dqn/gonso"
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
	n := gonso.New()
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

### How to get session token code

1. Select user and copy the link.

	![](docs/copy_link.png)

2. You can get session token code from query params.

	```
	npf71b963c1b7b6d119://auth#session_state=xxx&session_token_code=xxx...
	```
