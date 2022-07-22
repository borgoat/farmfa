.PHONY: test
test: test/go test/ffi

.PHONY: test/go
test/go:
	@go test ./...

.PHONY: test/ffi
test/ffi:
	@$(MAKE) -C ffitest test
