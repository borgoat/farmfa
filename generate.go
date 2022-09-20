package main

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest -o api/farmfa.gen.go -generate types,client,server,spec -package api api/farmfa.yaml
