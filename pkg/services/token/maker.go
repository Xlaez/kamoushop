package token

import (
	"time"
)

type Maker interface {
	CreateToken(user_id string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
