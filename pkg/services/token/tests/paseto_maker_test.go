package token_test

import (
	"kamoushop/pkg/services/token"
	"kamoushop/pkg/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := token.NewPasetoMaker(utils.RandomStr(32))
	require.NoError(t, err)

	user_id := primitive.NewObjectID()
	duration := time.Minute
	issuedAt := time.Now()
	expiresAt := time.Now().Add(duration)

	token, err := maker.CreateToken(user_id, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, user_id, payload.UserID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Second)
}
