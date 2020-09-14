package main

import (
	"fmt"
	"os"

	"github.com/dqn/gonso"
)

func run() error {
	sessiontToken, err := gonso.Login()
	if err != nil {
		return err
	}

	fmt.Println(sessiontToken)

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}
