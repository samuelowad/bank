package token

import (
	"fmt"
	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
	"time"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(key string) (Maker, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size, needs key size %d", chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(key),
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPlayLoad(username, duration)
	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

//	VerifyToken verifies the token is valid
func (maker *PasetoMaker) VerifyToken(token string) (*PayLoad, error) {
	payload := &PayLoad{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrorInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil

}
