package jwt

import (
	"log"
	"os"
	"time"
)

const (
	ValEmptyToken string = "Empty | NULL"
)

var (
	accessTokenSecret  []byte
	AccessTokenExpiry  time.Duration
	refreshTokenSecret []byte
	RefreshTokenExpiry time.Duration
)

func LoadConfig() {
	accessTokenSecret = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
	refreshTokenSecret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))

	var err error

	AccessTokenExpiry, err = time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRY"))
	if err != nil {
		log.Fatalln("[-] Error in JWT LoadConfig | parsing of Access Token Expiry.\n", err)
	}

	RefreshTokenExpiry, err = time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRY"))
	if err != nil {
		log.Fatalln("[-] Error in JWT LoadConfig | parsing of Refresh Token Expiry.\n", err)
	}

	log.Println("[+] JWT configuration load successful")
}
