package gwt

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// SecretKey — секретный ключ для подписи токенов
var SecretKey = []byte("NilfgaardAssirevarAnahid")

// Claims — структура, которая будет содержать информацию о пользователе
type Claims struct {
	jwt.RegisteredClaims
	UserAccess []string `json:"useraccess"`
	UserLogin  string   `json:"userlogin"`
}

// GenerateToken — временная функция для генерации JWT токена
func GenerateToken(userlogin string, useraccess []string) (string, error) {
	claims := Claims{
		UserLogin:  userlogin,
		UserAccess: useraccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Токен будет действителен 24 часа
			Issuer:    "AssirevarAnahid",                                  // Указание источника токена
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(SecretKey)
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}

	return signedToken, nil
}
