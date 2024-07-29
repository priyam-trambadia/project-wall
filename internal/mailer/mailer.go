package mailer

import (
	"context"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var (
	gmailService *gmail.Service
)

func SetupMailer() {
	ctx := context.Background()

	credentials, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalln("[-] Error in SetupMailer\n", err)
	}

	config, err := google.ConfigFromJSON(credentials, gmail.GmailSendScope)
	if err != nil {
		log.Fatalln("[-] Error in SetupMailer\n", err)
	}

	accessToken := os.Getenv("GMAIL_ACCESS_TOKEN")
	refreshToken := os.Getenv("GMAIL_REFRESH_TOKEN")

	token := oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	client := config.Client(ctx, &token)
	gmailService, err = gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalln("[-] Error in SetupMailer\n", err)
	}

	log.Println("[+] Mail service has been successfully set up.")
}
