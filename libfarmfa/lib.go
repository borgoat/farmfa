package main

import "C"

type ReturnCode int64

// TODO export
const (
	OK ReturnCode = iota
	ENOTARECIPIENT
	EINVALIDPLAYER
	EFAILEDTOCS
	EKEYGENFAIL
)

// An empty main function is needed to build a c-shared Go lib
func main() {}
