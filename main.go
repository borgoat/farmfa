package main

import "github.com/borgoat/farmfa/cmd"

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
