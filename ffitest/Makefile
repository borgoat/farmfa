.PHONY: all
all: clean build

.PHONY: test
test: lib_test
	@echo ${GODEBUG}
	@./lib_test

.PHONY: build
build: lib_test

.PHONY: clean
clean:
	@rm -f libfarmfa.h libfarmfa.so lib_test

lib_test: lib.c libfarmfa.so
	@$(CC) lib.c -Wl,-rpath=. -L. -lfarmfa -o lib_test

libfarmfa.so: export GODEBUG=cgocheck=2
libfarmfa.so:
	@go build -buildmode=c-shared -ldflags="-s -w" -o libfarmfa.so ../
