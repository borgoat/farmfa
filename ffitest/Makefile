lib.o: lib.c libfarmfa.so
	@$(CC) lib.c -Wl,-rpath=. -L. -lfarmfa -o test

libfarmfa.so:
	@go build -buildmode=c-shared -ldflags="-s -w" -o libfarmfa.so ../
