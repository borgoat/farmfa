package sessions

import (
	"github.com/giorgioazzinnaro/farmfa/shares"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestInMemory_Start(t *testing.T) {
	i := NewInMemory()

	session, err := i.Start(&mockShares[0])
	if err != nil {
		t.Fatalf("failed to start session")
	}

	assert.False(t, *session.Complete)
	assert.False(t, *session.Closed)
}

func TestInMemory_AddShare(t *testing.T) {
	i := NewInMemory()

	session, _ := i.Start(&mockShares[0])
	id := SessionIdentifier(*session.Id)

	i.AddShare(id, &mockShares[1])
	assert.False(t, *session.Complete)

	i.AddShare(id, &mockShares[2])
	assert.True(t, *session.Complete)
}

func TestInMemory_GenerateTOTP(t *testing.T) {

	i := NewInMemory()

	session, _ := i.Start(&mockShares[0])
	id := SessionIdentifier(*session.Id)

	i.AddShare(id, &mockShares[1])
	i.AddShare(id, &mockShares[2])
	assert.True(t, *session.Complete)

	code, err := i.GenerateTOTP(id)
	if err != nil {
		t.Fatalf("failed to generate shares: %v", err)
	}

	t.Logf("TOTP code: %s", code)
	assert.Regexp(t, regexp.MustCompile(`^\d{6}$`), code)
}

var mockShares = []shares.Token{
	{
		Secret:    "3RPPIDPM",
		Serial:    0,
		Threshold: 3,
		Total:     5,
		Share:     "tpMJGX0UeMX3834plGmbKsNWZ6bbl55MEtCz55zM8P8=RN9oMEaSzO1p6ZL2NRuNLZEc81IFuCauL3wgRQD9y2s=",
	},
	{
		Secret:    "3RPPIDPM",
		Serial:    1,
		Threshold: 3,
		Total:     5,
		Share:     "W2UOr4GzbOTECzQ560gy5bo6qMgIBNURY8QdRAInRYo=BObADeA-zJh2LT9g12XIwQDj0WV2vKmADg4eK1RWfdM=",
	},
	{
		Secret:    "3RPPIDPM",
		Serial:    2,
		Threshold: 3,
		Total:     5,
		Share:     "XsdVASqExS9Eq4QjQLmXFJvBJPrJO_N2I_f1HKu55io=emhByRKQL68nycgRrQPkv5gbTqFpN73KqatZRkD7ePg=",
	},
	{
		Secret:    "3RPPIDPM",
		Serial:    3,
		Threshold: 3,
		Total:     5,
		Share:     "bwzlCtaBzdfBDpqlhwLNSJ3tiNiuiBbcTp4MYjeXWZ0=skDfOWrdeJACUSQtsvnUt2t8GHwSqmZJfQSPvHXVP7s=",
	},
	{
		Secret:    "3RPPIDPM",
		Serial:    4,
		Threshold: 3,
		Total:     5,
		Share:     "Od6HvEHqI0OlETyjLCg47Kj0CFLQGZnMoRvd0aHuAew=38fpIDWiyO-rY62Rno4UCY5SEwXEKoPQjlFVi5RBBrU=",
	},
}
