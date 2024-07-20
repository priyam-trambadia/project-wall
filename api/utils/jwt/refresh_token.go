package jwt

// import (
// 	"errors"
// 	"time"

// 	"github.com/golang-jwt/jwt/v5"
// )

// func GenerateRefreshToken(userID int64) (string, error) {

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"user_id": userID,
// 		"expiry":  time.Now().Add(RefreshTokenExpiry).Unix(),
// 	})

// 	signedToken, err := token.SignedString(refreshTokenSecret)

// 	if err != nil {
// 		return "", err
// 	}

// 	return signedToken, nil
// }

// func ValidateRefreshToken(tokenString string) (int64, bool, error) {

// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return refreshTokenSecret, nil
// 	})

// 	if err != nil {
// 		return 0, false, err
// 	}

// 	if !token.Valid {
// 		return 0, false, nil
// 	}

// 	// Extract user ID from claims (if needed)
// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		return 0, false, errors.New("[-] Error in JWT ValidateRefreshToken | token claim type cast")
// 	}

// 	expiryf, ok1 := claims["expiry"].(float64)
// 	if !ok1 {
// 		return 0, false, errors.New("[-] Error in JWT ValidateRefreshToken | token claim->expiry type cast")

// 	}
// 	expiry := int64(expiryf)
// 	if time.Now().Unix() > expiry {
// 		return 0, false, nil
// 	}

// 	userIDf, ok2 := claims["user_id"].(float64)
// 	if !ok2 {
// 		return 0, false, errors.New("[-] Error in JWT ValidateRefreshToken | token claim->user_id type cast")
// 	}
// 	userID := int64(userIDf)

// 	return userID, true, nil
// }
