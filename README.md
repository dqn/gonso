# go-nso

Nintendo Switch Online API Wrapper

## Installation

```bash
$ go get github.com/dqn/go-nso
```

## Usage

```go
package main

import "github.com/dqn/go-nso"

func main() {
	n := nso.New()
	n.Auth()
}

```

```bash
authorize by visiting this url: https://accounts.nintendo.com/connect/1.0.0/authorize?xxx=xxx
session token code: <input your session token code>
```

Credentials are cached as `./nso.json`
