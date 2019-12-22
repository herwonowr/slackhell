package main

import (
	"os"

	"github.com/herwonowr/slackhell/cmd/app"
)

func main() {
	if err := app.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
