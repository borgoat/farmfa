package main

import "github.com/giorgioazzinnaro/farmfa/cmd"

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
