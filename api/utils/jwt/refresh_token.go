package jwt

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateRefreshToken(userID int64) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"expiry":  time.Now().Add(refreshTokenExpiry).Unix(),
	})

	signedToken, err := token.SignedString(refreshTokenSecret)

	if err != nil {
		log.Fatalln("[-] Error in JWT generate refresh token signing")
	}

	return signedToken
}

func ValidateRefreshToken(tokenString string) (int64, bool) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return refreshTokenSecret, nil
	})

	if err != nil {
		log.Fatalln("[-] Error in JWT validate refresh token parsing")
	}

	if !token.Valid {
		return 0, false
	}

	// Extract user ID from claims (if needed)
	claims := token.Claims.(jwt.MapClaims)

	expiryf := claims["expiry"].(float64)
	expiry := int64(expiryf)
	if time.Now().Unix() > expiry {
		return 0, false
	}

	userIDf, _ := claims["user_id"].(float64)
	userID := int64(userIDf)

	return userID, true
}
