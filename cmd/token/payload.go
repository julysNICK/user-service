package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type Payload struct {
	ID       uuid.UUID `json:"id"`
	UserName string    `json:"user_name"`
	Expires  time.Time `json:"expires"`
}

func NewPayload(userName string) (*Payload, error) {
	token, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:       token,
		UserName: userName,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
	}

	return payload, nil
}

func (p *Payload) Valid() error {

	if p.ID == uuid.Nil {
		return ErrInvalidToken
	}

	if time.Now().After(p.Expires) {
		return ErrExpiredToken
	}

	return nil

}
