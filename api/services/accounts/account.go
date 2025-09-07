package accounts

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Account struct {
	ID           string `bson:"_id"`
	Email        string `bson:"email"`
	RefreshToken string `bson:"refresh_token"`
}

var (
	_account_prefix = "ACC"
	secretKey       = "secret"

	ErrMissingEmail = errors.New("missing email")
)

func NewAccount(email string) (*Account, error) {
	if email == "" {
		return nil, ErrMissingEmail
	}

	return &Account{
		ID:    _account_prefix + strings.ReplaceAll(uuid.New().String(), "-", ""),
		Email: email,
	}, nil
}

func (a *Account) GenerateTokens() (string, string, error) {
	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    a.ID,
		"email": a.Email,
		"exp":   jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		"iat":   jwt.NewNumericDate(time.Now()),
	})

	accessToken, err := accessClaims.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    a.ID,
		"email": a.Email,
		"exp":   jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		"iat":   jwt.NewNumericDate(time.Now()),
	})

	refreshToken, err := refreshClaims.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	a.RefreshToken = refreshToken

	return accessToken, refreshToken, nil

}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
