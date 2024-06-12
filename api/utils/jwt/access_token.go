package jwt

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID int64) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"expiry":  time.Now().Add(accessTokenExpiry).Unix(),
	})

	signedToken, err := token.SignedString(accessTokenSecret)

	if err != nil {
		log.Fatalln("[-] Error in JWT generate access token signing")
	}

	return signedToken
}

func ValidateAccessToken(tokenString string) (int64, bool) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return accessTokenSecret, nil
	})

	if err != nil {
		log.Fatalln("[-] Error in JWT validate access token parsing")
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
