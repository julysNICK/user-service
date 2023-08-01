package token

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

type JWTMaker struct {
	secretkey string
}

func NewJWTMaker(secretkey string) (Maker, error) {
	if len(secretkey) < 10 {
		return nil, ErrInvalidToken
	}

	return &JWTMaker{secretkey}, nil
}

func (maker *JWTMaker) CreateToken(username string) (string, *Payload, error) {
	payload, err := NewPayload(username)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	token, err := jwtToken.SignedString([]byte(maker.secretkey))

	if err != nil {
		return "", nil, err
	}

	return token, payload, nil

}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(maker.secretkey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		verr, ok := err.(*jwt.ValidationError)

		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}

		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)

	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil

}
