package jwt

import (
	"log"
	"os"
	"time"
)

var (
	accessTokenSecret  []byte
	accessTokenExpiry  time.Duration
	refreshTokenSecret []byte
	refreshTokenExpiry time.Duration
)

func LoadConfig() {
	accessTokenSecret = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
	refreshTokenSecret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))

	var err error

	accessTokenExpiry, err = time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRY"))
	if err != nil {
		log.Fatalln("[-] JWT error in parsing of Access Token Expiry")
	}

	refreshTokenExpiry, err = time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRY"))
	if err != nil {
		log.Fatalln("[-] JWT error in parsing of Refresh Token Expiry")
	}

	log.Println("[+] JWT configuration load successful")
}
