# gonso

Nintendo Switch Online API wrapper written in Go.

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

func main() {
  sessionToken, err := gonso.Login()
  if err != nil {
    // Handle error.
  }

  // if you save the session token, you can skip the login next time
  // err := ioutil.WriteFile("session_token.txt", []byte(sessiontToken), 0644)
  // if err != nil {
  //   // Handle error.
  // }

  accessToken, err := gonso.Auth(sessionToken)
  if err != nil {
    // Handle error.
  }
}
```

```bash
authenticate by visiting this url: https://accounts.nintendo.com/connect/1.0.0/authorize?...
session token code: <enter-your-session-token-code>
```

### How to get session token code

1.  Visit the generated URL, select user and copy the link.

    ![](docs/copy_link.png)

2.  You can get session token code from the URL.

    ```
    npf71b963c1b7b6d119://auth#...&session_token_code=<session-token-code>...
    ```

## See Also

This project is inspired by [splatnet2statink](https://github.com/frozenpandaman/splatnet2statink).

## License

MIT
