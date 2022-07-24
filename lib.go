package main

/*
#include <stdint.h> // for uintptr_t
#include <stdlib.h>
#include <string.h>

typedef uintptr_t fm_dealer_t;

typedef struct fm_keypair {
	char *public_key;
	char *private_key;
} fm_keypair;

typedef struct fm_encrypted_toc {
	char *recipient;
	char *encrypted_toc;
} fm_encrypted_toc;

typedef struct fm_encrypted_tocs {
	fm_encrypted_toc *items;
	size_t length;
} fm_encrypted_tocs;

*/
import "C"

import (
	"filippo.io/age"
	"github.com/borgoat/farmfa/deal"
	"runtime/cgo"
	"unsafe"
)

type dealerCtx struct {
	players []*deal.Player
	secret  string
	note    string
	tocs    map[string]string
	tocsIdx int32
}

//export fm_dealer_init
func fm_dealer_init(handle *C.fm_dealer_t) int32 {
	*handle = C.fm_dealer_t(cgo.NewHandle(&dealerCtx{}))
	return 0
}

//export fm_dealer_free
func fm_dealer_free(handle C.fm_dealer_t) {
	cgo.Handle(handle).Delete()
}

//export fm_dealer_add_player
func fm_dealer_add_player(handle C.fm_dealer_t, recipient, key *C.char) int32 {
	ctx := dealerContextFromHandle(handle)
	r := C.GoString(recipient)
	k := C.GoString(key)

	ageRecipient, err := age.ParseX25519Recipient(k)
	if err != nil {
		return 1
	}
	player, err := deal.NewPlayer(r, deal.EncryptWithAge(ageRecipient))
	if err != nil {
		return 2
	}

	ctx.players = append(ctx.players, player)

	return 0
}

//export fm_dealer_set_secret
func fm_dealer_set_secret(handle C.fm_dealer_t, secret *C.char) int32 {
	ctx := dealerContextFromHandle(handle)
	s := C.GoString(secret)

	ctx.secret = s

	return 0
}

//export fm_dealer_set_note
func fm_dealer_set_note(handle C.fm_dealer_t, note *C.char) int32 {
	ctx := dealerContextFromHandle(handle)
	n := C.GoString(note)

	ctx.note = n

	return 0
}

//export fm_dealer_create_tocs
func fm_dealer_create_tocs(handle C.fm_dealer_t, encrypted_tocs *C.fm_encrypted_tocs) int32 {
	ctx := dealerContextFromHandle(handle)

	tocs, err := deal.CreateTocs(ctx.note, ctx.secret, ctx.players, 3)
	if err != nil {
		return 1
	}

	(*encrypted_tocs).length = C.size_t(len(tocs))
	(*encrypted_tocs).items = (*C.fm_encrypted_toc)(C.calloc(C.size_t(len(tocs)), C.size_t(unsafe.Sizeof(C.fm_encrypted_toc{}))))

	var tocsSlice = make([]C.fm_encrypted_toc, len(tocs))

	i := 0
	for r, t := range tocs {
		tocsSlice[i] = C.fm_encrypted_toc{
			recipient:     C.CString(r),
			encrypted_toc: C.CString(t),
		}
		i++
	}

	C.memcpy(unsafe.Pointer((*encrypted_tocs).items), unsafe.Pointer(&tocsSlice[0]), C.size_t(uintptr(len(tocsSlice))*unsafe.Sizeof(C.fm_encrypted_toc{})))

	return 0
}

func dealerContextFromHandle(handle C.fm_dealer_t) *dealerCtx {
	return cgo.Handle(handle).Value().(*dealerCtx)
}

//export fm_player_create_key
func fm_player_create_key(keypair *C.fm_keypair) int32 {
	id, err := age.GenerateX25519Identity()
	if err != nil {
		return 1
	}

	*keypair = C.fm_keypair{
		public_key:  C.CString(id.Recipient().String()),
		private_key: C.CString(id.String()),
	}

	return 0
}

//export fm_player_keypair_free
func fm_player_keypair_free(keypair *C.fm_keypair) int32 {
	C.free(unsafe.Pointer(keypair.public_key))
	C.free(unsafe.Pointer(keypair.private_key))
	return 0
}
