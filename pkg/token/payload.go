package token

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrorExpiredToken = errors.New("token has expired")
	ErrorInvalidToken = errors.New("token is invalid")
)

type PayLoad struct {
	ID        uuid.UUID `json:"json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func NewPlayLoad(username string, duration time.Duration) (*PayLoad, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &PayLoad{
		ID:        id,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *PayLoad) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrorExpiredToken
	}
	return nil
}
