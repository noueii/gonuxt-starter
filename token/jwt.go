package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

var ErrInvalidJWTToken = errors.New("invalid jwt token")

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key length")
	}

	return &JWTMaker{
		secretKey: secretKey,
	}, nil
}

func (m *JWTMaker) CreateToken(userID uuid.UUID, email string, username string, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, email, username, role, duration)

	if err != nil {
		return "", nil, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	signedToken, err := jwtToken.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", nil, err
	}

	return signedToken, payload, nil
}

func (m *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}

		return []byte(m.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidJWTToken
	}

	return payload, nil
}
