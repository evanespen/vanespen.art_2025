package security

import (
	"errors"
	"github.com/evanespen/vanespen.art_2025/configs"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func Authenticate(password string) (string, error) {
	err := bcrypt.CompareHashAndPassword([]byte(configs.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(configs.Salt))
	if err != nil {
		return "", errors.New("unable to sign token")
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(configs.Salt), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok {
		return nil
	} else {
		return err
	}
}
