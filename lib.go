package main

/*
#include <stdint.h> // for uintptr_t

typedef uintptr_t fm_dealer_t;
*/
import "C"

import (
	"filippo.io/age"
	"github.com/borgoat/farmfa/deal"
	"runtime/cgo"
)

type dealerCtx struct {
	players []*deal.Player
	secret  string
	note    string
}

//export fm_dealer_init
func fm_dealer_init() C.fm_dealer_t {
	return C.uintptr_t(cgo.NewHandle(&dealerCtx{}))
}

//export fm_dealer_free
func fm_dealer_free(handle C.fm_dealer_t) {
	cgo.Handle(handle).Delete()
}

//export fm_dealer_add_player
func fm_dealer_add_player(handle C.fm_dealer_t, recipient, key *C.char) int64 {
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
func fm_dealer_set_secret(handle C.fm_dealer_t, secret *C.char) int64 {
	ctx := dealerContextFromHandle(handle)
	s := C.GoString(secret)

	ctx.secret = s

	return 0
}

//export fm_dealer_set_note
func fm_dealer_set_note(handle C.fm_dealer_t, note *C.char) int64 {
	ctx := dealerContextFromHandle(handle)
	n := C.GoString(note)

	ctx.note = n

	return 0
}

//export fm_dealer_create_tocs
func fm_dealer_create_tocs(handle C.fm_dealer_t) int64 {
	ctx := dealerContextFromHandle(handle)

	_, err := deal.CreateTocs("", ctx.secret, ctx.players, 3)
	if err != nil {
		return 1
	}

	return 0
}

func dealerContextFromHandle(handle C.fm_dealer_t) *dealerCtx {
	return cgo.Handle(handle).Value().(*dealerCtx)
}

//export fm_player_create_key
func fm_player_create_key(public_key_buffer *C.char, private_key_buffer *C.char) int64 {
	id, err := age.GenerateX25519Identity()
	if err != nil {
		return 1
	}

	*public_key_buffer = *C.CString(id.Recipient().String())
	*private_key_buffer = *C.CString(id.String())

	return 0
}
