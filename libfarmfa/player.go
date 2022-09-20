package main

/*
#include <stdint.h> // for uintptr_t
#include <stdlib.h>
#include <string.h>

typedef uintptr_t fm_player_t;

typedef struct fm_keypair {
	char *public_key;
	char *private_key;
} fm_keypair;
*/
import "C"

import (
	"filippo.io/age"
	"filippo.io/age/armor"
	"io"
	"runtime/cgo"
	"strings"
	"unsafe"
)

type playerCtx struct {
	identities []age.Identity
}

func playerContextFromHandle(handle C.fm_player_t) *playerCtx {
	return cgo.Handle(handle).Value().(*playerCtx)
}

//export fm_player_create_key
func fm_player_create_key(keypair *C.fm_keypair) ReturnCode {
	id, err := age.GenerateX25519Identity()
	if err != nil {
		return EKEYGENFAIL
	}

	*keypair = C.fm_keypair{
		public_key:  C.CString(id.Recipient().String()),
		private_key: C.CString(id.String()),
	}

	return OK
}

//export fm_player_keypair_free
func fm_player_keypair_free(keypair *C.fm_keypair) ReturnCode {
	C.free(unsafe.Pointer(keypair.public_key))
	C.free(unsafe.Pointer(keypair.private_key))
	return OK
}

//export fm_player_init
func fm_player_init(handle *C.fm_player_t) ReturnCode {
	*handle = C.fm_player_t(cgo.NewHandle(&playerCtx{}))
	return OK
}

//export fm_player_free
func fm_player_free(handle C.fm_player_t) {
	cgo.Handle(handle).Delete()
}

//export fm_player_load_identity
func fm_player_load_identity(handle C.fm_player_t, private_key *C.char) ReturnCode {
	ctx := playerContextFromHandle(handle)

	id, err := age.ParseX25519Identity(C.GoString(private_key))
	if err != nil {
		return 1
	}

	ctx.identities = append(ctx.identities, id)
	return OK
}

//export fm_player_decrypt
func fm_player_decrypt(handle C.fm_player_t, armored *C.char, decrypted **C.char) ReturnCode {
	ctx := playerContextFromHandle(handle)

	s := C.GoString(armored)
	r := armor.NewReader(strings.NewReader(s))
	o, err := age.Decrypt(r, ctx.identities...)
	if err != nil {
		return 1
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, o)
	if err != nil {
		return 2
	}

	*decrypted = C.CString(buf.String())
	return OK
}
