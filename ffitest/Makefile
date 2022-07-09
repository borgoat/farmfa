.PHONY: all clean build test

all: clean build

test: lib_test
	@./lib_test

build: lib_test

clean:
	@rm -f libfarmfa.h libfarmfa.so lib_test

lib_test: lib.c libfarmfa.so
	@$(CC) lib.c -Wl,-rpath=. -L. -lfarmfa -o lib_test

libfarmfa.so:
	@go build -buildmode=c-shared -ldflags="-s -w" -o libfarmfa.so ../
