package jwtmiddleware

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ibamaulana/golang-kumparan/config"
	"golang.org/x/crypto/bcrypt"
)

type TokenRequestStructure struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
}

type TokenResponse struct {
	AccessToken string  `json:"access_token"`
	TokenType   string  `json:"token_type"`
	ExpiredIn   float64 `json:"expired_in"`
	ExpiredAt   int64   `json:"expired_at"`
}

type customAuth struct {
	signature []byte
}

type CustomAuth interface {
	GenerateToken(data TokenRequestStructure) (response TokenResponse, err error)
}

func NewCustomAuth(signature []byte) CustomAuth {
	return &customAuth{signature}
}

func (cAuth *customAuth) GenerateToken(data TokenRequestStructure) (response TokenResponse, err error) {
	cfg := config.NewConfig()

	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)

	expiredIn := time.Hour * (24 * 7)
	expiredAt := time.Now().Add(time.Hour * (24 * 7))

	myCrypt, err := bcrypt.GenerateFromPassword([]byte(cfg.GetString("app.signature")), 8)
	if err != nil {
		return
	}

	claims["user_id"] = data.UserID
	claims["email"] = data.Email
	claims["hash"] = string(myCrypt)
	claims["exp"] = expiredIn

	tokenString, _ := token.SignedString(cAuth.signature)

	response.AccessToken = tokenString
	response.TokenType = "Bearer"
	response.ExpiredAt = expiredAt.Unix()
	response.ExpiredIn = expiredIn.Seconds()

	return
}
