package main

import (
	"os"

	"github.com/borgoat/farmfa/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
