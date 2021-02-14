package session_test

import (
	"bytes"
	"filippo.io/age"
	"filippo.io/age/armor"
	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/giorgioazzinnaro/farmfa/session"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestInMemory_CreateSession(t *testing.T) {
	m := session.NewInMemory()
	creds, err := m.CreateSession(&api.Toc{
		GroupId:        "GIWYCRKS",
		GroupSize:      20,
		GroupThreshold: 10,
		Note:           nil,
		Share:          "EjidONG9xkvJAzhZzn3bKua4_5rChSfQvbc3q2sFqQ0=mL_35wzKgcTUF4iOiIdI5iYy8eXCjDdCI-4XibSVqWo=",
		TocId:          "RSIFMOCY",
	})
	assert.NoError(t, err)
	assert.Equal(t, "pending", creds.Status)

	// assert that TEK is a valid age recipient
	_, err = age.ParseX25519Recipient(creds.Tek)
	assert.NoError(t, err)
}

func TestInMemory_AddToc_valid(t *testing.T) {
	m := session.NewInMemory()
	creds, err := m.CreateSession(&api.Toc{
		GroupId:        "4TRP4K4R",
		GroupSize:      20,
		GroupThreshold: 10,
		Note:           nil,
		Share:          "WLkXB9ZOScDGlzuaVDYF61UkH68In_lrQ1WoZSUN53I=B4D_B-HZoDuGvRWzi9KOdSFrecpjIVqsjLs-MCtNlP0=",
		TocId:          "LLOAXTBV",
	})
	assert.NoError(t, err)

	err = m.AddToc(creds.Id, "-----BEGIN AGE ENCRYPTED FILE-----\nYWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSBrVHUyVVQ2UHQzdHk2RGYw\nY2xGMWwvNUwvNGRwaThGUHk5NU9LZmtDdndVClg1S1ZqdGtpMVdWbThjem5XY09x\naFpCT094ZGZTSVBuTHpxSUt0NC9HUjgKLS0tIHRnVGFmcDJ5ZnQyMDdhNm9ocFVG\ndVQxdmJrL2FuZmZ4OHZ3MVdGV1NtZXcKfLvPWhP8je+azJn3hwb/QeQ4lV91rDca\niuVX2+ch9Vks5/mKx6hf0HhDs2Ak7gifdfJAuzzyPp2ap+Oy5rQleIUf7lCmPCq5\nY0Ued1ohwpeMMa7gMFL5cjOrwGAHDqZ4ur9xk1uKfS7wTJ3fPp/xPJPJnAmOT8Xg\nMDCcObCFaY/5ewWWPJHqVmt+MhmNmjMrO5wIzK7qGnQbNRUAAjalWxjjvC+V//SO\n+PjDUNVDBnVJk0kZE/GZh5YcHuVC6poPgZyPswsd5jF/P5d9zD7b2+rJujyrJTus\nKMnCDLMD0ut6YcPNP/fulQ==\n-----END AGE ENCRYPTED FILE-----\n")
	assert.NoError(t, err)
}

func TestInMemory_AddToc_empty(t *testing.T) {
	m := session.NewInMemory()
	creds, err := m.CreateSession(&api.Toc{
		GroupId:        "GIWYCRKS",
		GroupSize:      20,
		GroupThreshold: 10,
		Note:           nil,
		Share:          "EjidONG9xkvJAzhZzn3bKua4_5rChSfQvbc3q2sFqQ0=mL_35wzKgcTUF4iOiIdI5iYy8eXCjDdCI-4XibSVqWo=",
		TocId:          "RSIFMOCY",
	})
	assert.NoError(t, err)

	err = m.AddToc(creds.Id, "")
	assert.ErrorIs(t, err, session.ErrEmptyToc)
}

func TestInMemory_AddToc_notEncrypted(t *testing.T) {
	m := session.NewInMemory()
	creds, err := m.CreateSession(&api.Toc{
		GroupId:        "GIWYCRKS",
		GroupSize:      20,
		GroupThreshold: 10,
		Note:           nil,
		Share:          "EjidONG9xkvJAzhZzn3bKua4_5rChSfQvbc3q2sFqQ0=mL_35wzKgcTUF4iOiIdI5iYy8eXCjDdCI-4XibSVqWo=",
		TocId:          "RSIFMOCY",
	})
	assert.NoError(t, err)

	err = m.AddToc(creds.Id, "this is not an age armored string")
	assert.ErrorIs(t, err, session.ErrTocIsNotEncrypted)
}

func TestInMemory_AddToc_alreadyExists(t *testing.T) {
	m := session.NewInMemory()
	creds, err := m.CreateSession(&api.Toc{
		GroupId:        "GIWYCRKS",
		GroupSize:      20,
		GroupThreshold: 10,
		Note:           nil,
		Share:          "EjidONG9xkvJAzhZzn3bKua4_5rChSfQvbc3q2sFqQ0=mL_35wzKgcTUF4iOiIdI5iYy8eXCjDdCI-4XibSVqWo=",
		TocId:          "RSIFMOCY",
	})
	assert.NoError(t, err)

	err = m.AddToc(creds.Id, "-----BEGIN AGE ENCRYPTED FILE-----\nYWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSBrVHUyVVQ2UHQzdHk2RGYw\nY2xGMWwvNUwvNGRwaThGUHk5NU9LZmtDdndVClg1S1ZqdGtpMVdWbThjem5XY09x\naFpCT094ZGZTSVBuTHpxSUt0NC9HUjgKLS0tIHRnVGFmcDJ5ZnQyMDdhNm9ocFVG\ndVQxdmJrL2FuZmZ4OHZ3MVdGV1NtZXcKfLvPWhP8je+azJn3hwb/QeQ4lV91rDca\niuVX2+ch9Vks5/mKx6hf0HhDs2Ak7gifdfJAuzzyPp2ap+Oy5rQleIUf7lCmPCq5\nY0Ued1ohwpeMMa7gMFL5cjOrwGAHDqZ4ur9xk1uKfS7wTJ3fPp/xPJPJnAmOT8Xg\nMDCcObCFaY/5ewWWPJHqVmt+MhmNmjMrO5wIzK7qGnQbNRUAAjalWxjjvC+V//SO\n+PjDUNVDBnVJk0kZE/GZh5YcHuVC6poPgZyPswsd5jF/P5d9zD7b2+rJujyrJTus\nKMnCDLMD0ut6YcPNP/fulQ==\n-----END AGE ENCRYPTED FILE-----\n")
	assert.NoError(t, err)

	err = m.AddToc(creds.Id, "-----BEGIN AGE ENCRYPTED FILE-----\nYWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSBrVHUyVVQ2UHQzdHk2RGYw\nY2xGMWwvNUwvNGRwaThGUHk5NU9LZmtDdndVClg1S1ZqdGtpMVdWbThjem5XY09x\naFpCT094ZGZTSVBuTHpxSUt0NC9HUjgKLS0tIHRnVGFmcDJ5ZnQyMDdhNm9ocFVG\ndVQxdmJrL2FuZmZ4OHZ3MVdGV1NtZXcKfLvPWhP8je+azJn3hwb/QeQ4lV91rDca\niuVX2+ch9Vks5/mKx6hf0HhDs2Ak7gifdfJAuzzyPp2ap+Oy5rQleIUf7lCmPCq5\nY0Ued1ohwpeMMa7gMFL5cjOrwGAHDqZ4ur9xk1uKfS7wTJ3fPp/xPJPJnAmOT8Xg\nMDCcObCFaY/5ewWWPJHqVmt+MhmNmjMrO5wIzK7qGnQbNRUAAjalWxjjvC+V//SO\n+PjDUNVDBnVJk0kZE/GZh5YcHuVC6poPgZyPswsd5jF/P5d9zD7b2+rJujyrJTus\nKMnCDLMD0ut6YcPNP/fulQ==\n-----END AGE ENCRYPTED FILE-----\n")
	assert.ErrorIs(t, err, session.ErrTocAlreadyExists)
}

func TestInMemory_DecryptTocs(t *testing.T) {
	m := session.NewInMemory()
	creds, err := m.CreateSession(&api.Toc{
		GroupId:        "J7UHQPZK",
		GroupSize:      5,
		GroupThreshold: 2,
		Note:           nil,
		Share:          "5Ovpu-PKEeYXx5ebiQhzU_AT0Z79POf8GGkskDp3its=urkBkVXr-pYjIvTt1ch2YJILCScAoRquLoX_VBxxps4=",
		TocId:          "TFW52GAK",
	})
	assert.NoError(t, err)

	tek, err := age.ParseX25519Recipient(creds.Tek)
	assert.NoError(t, err)

	var out bytes.Buffer
	aw := armor.NewWriter(&out)
	w, err := age.Encrypt(aw, tek)
	assert.NoError(t, err)

	_, err = io.WriteString(w, `{"group_id":"J7UHQPZK","group_size":5,"group_threshold":2,"note":"test_basic","share":"zxRrozuUaCMgn_u6ajZStlV7RKwhp0keT9aQoXAEruI=nfx2CPJfKiFM32zLmtxHjV94OlZOgBevV1Whrx-lslU=","toc_id":"K5FSSJSV"}`)
	assert.NoError(t, err)

	assert.NoError(t, w.Close())
	assert.NoError(t, aw.Close())

	err = m.AddToc(creds.Id, out.String())

	tocs, err := m.DecryptTocs(creds.Id, &creds.SessionKeyEncryptionKey)
	assert.NoError(t, err)
	assert.Len(t, tocs, 2)
}
