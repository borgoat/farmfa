package main

import "github.com/giorgioazzinnaro/multi-farmer-authentication/cmd"

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
