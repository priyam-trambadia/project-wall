package jwt

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/priyam-trambadia/project-wall/internal/logger"
)

const (
	ValEmptyToken string = "Empty | NULL"
)

var (
	accessTokenSecret  []byte
	AccessTokenExpiry  time.Duration
	refreshTokenSecret []byte
	RefreshTokenExpiry time.Duration

	userActivationTokenSecret    []byte
	userActivationTokenExpiry    time.Duration
	userPasswordResetTokenSecret []byte
	userPasswordResetTokenExpiy  time.Duration
)

func LoadConfig() {
	accessTokenSecret = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
	refreshTokenSecret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))
	userActivationTokenSecret = []byte(os.Getenv("USER_ACTIVATION_TOKEN_SECRET"))
	userPasswordResetTokenSecret = []byte(os.Getenv("USER_PASSWORD_RESET_TOKEN_SECRET"))

	var err error

	AccessTokenExpiry, err = time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRY"))
	if err != nil {
		log.Fatalln("[-] Error in JWT LoadConfig | parsing of Access Token Expiry.\n", err)
	}

	RefreshTokenExpiry, err = time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRY"))
	if err != nil {
		log.Fatalln("[-] Error in JWT LoadConfig | parsing of Refresh Token Expiry.\n", err)
	}

	userActivationTokenExpiry, err = time.ParseDuration(os.Getenv("USER_ACTIVATION_TOKEN_EXPIRY"))
	if err != nil {
		log.Fatalln("[-] Error in JWT LoadConfig | parsing of User  Activation Token Expiry.\n", err)
	}

	userPasswordResetTokenExpiy, err = time.ParseDuration(os.Getenv("USER_PASSWORD_RESET_TOKEN_EXPIRY"))
	if err != nil {
		log.Fatalln("[-] Error in JWT LoadConfig | parsing of User Password Reset Token Expiry.\n", err)
	}

	log.Println("[+] JWT configuration load successful")
}

func generateToken(userID int64, tokenExpiry time.Duration, tokenSecret []byte) (string, error) {
	logger := logger.Logger{Caller: "generateToken utils/jwt"}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"expiry":  time.Now().Add(tokenExpiry).Unix(),
	})

	signedToken, err := token.SignedString(tokenSecret)
	if err != nil {
		return "", logger.AppendError(err)
	}

	return signedToken, nil
}

func validateToken(tokenString string, tokenSecret []byte) (int64, bool, error) {
	logger := logger.Logger{Caller: "validateToken utils/jwt"}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return tokenSecret, nil
	})

	if err != nil {
		return 0, false, logger.AppendError(err)
	}

	if !token.Valid {
		return 0, false, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err := logger.AppendError(errors.New("error in token claim type cast"))
		return 0, false, err
	}

	expiryf, ok := claims["expiry"].(float64)
	if !ok {
		err := logger.AppendError(errors.New("error in token claim->expiry type cast"))
		return 0, false, err
	}

	expiry := int64(expiryf)
	if time.Now().Unix() > expiry {
		return 0, false, nil
	}

	userIDf, ok := claims["user_id"].(float64)
	if !ok {
		err := logger.AppendError(errors.New("error in token claim->user_id type cast"))
		return 0, false, err
	}
	userID := int64(userIDf)

	return userID, true, nil
}
