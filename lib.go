package main

/*
#include <stdint.h> // for uintptr_t
#include <stdlib.h>
#include <string.h>

typedef uintptr_t fm_dealer_t;
typedef uintptr_t fm_player_t;

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
	"filippo.io/age/armor"
	"github.com/borgoat/farmfa/deal"
	"github.com/hashicorp/go-multierror"
	"io"
	"runtime/cgo"
	"strings"
	"unsafe"
)

type dealerCtx struct {
	players []*deal.Player
	secret  string
	note    string
	tocs    map[string]string
	tocsIdx ReturnCode
	err     error
}

func dealerContextFromHandle(handle C.fm_dealer_t) *dealerCtx {
	return cgo.Handle(handle).Value().(*dealerCtx)
}

type ReturnCode int64

// export
const (
	OK ReturnCode = iota
	ENOTARECIPIENT
	EINVALIDPLAYER
	EFAILEDTOCS
	EKEYGENFAIL
)

//export fm_dealer_init
func fm_dealer_init(handle *C.fm_dealer_t) ReturnCode {
	*handle = C.fm_dealer_t(cgo.NewHandle(&dealerCtx{}))
	return OK
}

//export fm_dealer_free
func fm_dealer_free(handle C.fm_dealer_t) {
	cgo.Handle(handle).Delete()
}

//export fm_dealer_add_player
func fm_dealer_add_player(handle C.fm_dealer_t, recipient, key *C.char) ReturnCode {
	ctx := dealerContextFromHandle(handle)
	r := C.GoString(recipient)
	k := C.GoString(key)

	ageRecipient, err := age.ParseX25519Recipient(k)
	if err != nil {
		ctx.err = multierror.Append(ctx.err, err)
		return ENOTARECIPIENT
	}
	player, err := deal.NewPlayer(r, deal.EncryptWithAge(ageRecipient))
	if err != nil {
		ctx.err = multierror.Append(ctx.err, err)
		return EINVALIDPLAYER
	}

	ctx.players = append(ctx.players, player)

	return OK
}

//export fm_dealer_set_secret
func fm_dealer_set_secret(handle C.fm_dealer_t, secret *C.char) ReturnCode {
	ctx := dealerContextFromHandle(handle)
	s := C.GoString(secret)

	ctx.secret = s

	return OK
}

//export fm_dealer_set_note
func fm_dealer_set_note(handle C.fm_dealer_t, note *C.char) ReturnCode {
	ctx := dealerContextFromHandle(handle)
	n := C.GoString(note)

	ctx.note = n

	return OK
}

//export fm_dealer_create_tocs
func fm_dealer_create_tocs(handle C.fm_dealer_t, encrypted_tocs *C.fm_encrypted_tocs) ReturnCode {
	ctx := dealerContextFromHandle(handle)

	tocs, err := deal.CreateTocs(ctx.note, ctx.secret, ctx.players, 3)
	if err != nil {
		ctx.err = multierror.Append(ctx.err, err)
		return EFAILEDTOCS
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

	return OK
}

//export fm_dealer_get_errors
func fm_dealer_get_errors(handle C.fm_dealer_t, errors **C.char) {
	ctx := dealerContextFromHandle(handle)

	*errors = C.CString(ctx.err.Error())
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

type playerCtx struct {
	identities []age.Identity
}

func playerContextFromHandle(handle C.fm_player_t) *playerCtx {
	return cgo.Handle(handle).Value().(*playerCtx)
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
