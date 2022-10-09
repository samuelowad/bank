package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const minSecretsKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {

	if len(secretKey) < minSecretsKeySize {
		return nil, fmt.Errorf("secret key must be at least %d characters", minSecretsKeySize)
	}

	return &JWTMaker{secretKey: secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPlayLoad(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

//	VerifyToken verifies the token is valid
func (maker *JWTMaker) VerifyToken(token string) (*PayLoad, error) {
	keyFun := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrorInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &PayLoad{}, keyFun)
	if err != nil {
		vErr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(vErr.Inner, ErrorExpiredToken) {
			return nil, ErrorExpiredToken
		}
		return nil, ErrorInvalidToken
	}
	payload, ok := jwtToken.Claims.(*PayLoad)
	if !ok {
		return nil, ErrorInvalidToken
	}
	return payload, nil
}
