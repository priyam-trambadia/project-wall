package mailer

import (
	"encoding/base64"
	"fmt"

	"github.com/priyam-trambadia/project-wall/internal/logger"
	"google.golang.org/api/gmail/v1"
)

func send(mail string) error {
	logger := logger.Logger{Caller: "send mail"}

	messageStr := []byte(mail)

	var message gmail.Message
	message.Raw = base64.URLEncoding.EncodeToString(messageStr)

	_, err := gmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		err = logger.AppendError(err)
	}

	return err
}

func mailHead(to string, subject string) string {
	return fmt.Sprintf(
		"From: ProjectWall <priyam.projectwall@gmail.com>\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"utf-8\"\r\n\r\n",
		to,
		subject,
	)
}

func mailBody(body string) string {
	return fmt.Sprintf(`
		<!doctype html>
		<html>
		<head>
				<meta name="viewpoint" content="width=device-width"/>
				<meta http-equiv="Content-Type" content="text/html"; charset="UTF-8"/>
		</head>

		<body>
			%s
		</body>`,
		body,
	)
}

func SendUserActivationMail(to, link string) error {
	logger := logger.Logger{Caller: "SendUserActivationMail mail"}

	subject := "Activation Email"
	body := `
		<p>To finish setting up your account, we just need to make sure this email address is yours.</p>
		<p>Click the link below to activate your account:
			<a href=%s>Activate</a>
		</p>
		<p>Or use this: <br>%s</p>
		<p>If you didn't request this, you can safely ignore this email.</p>
		<br/>
		<p>Thanks,</p>
		<p>ProjectWall team</p>
	`

	body = fmt.Sprintf(body, link, link)
	body = mailHead(to, subject) + mailBody(body)

	err := send(body)
	if err != nil {
		err = logger.AppendError(err)
	}

	return err
}

func SendUserPasswordResetMail(to, link string) error {
	logger := logger.Logger{Caller: "SendUserPasswordResetMail mail"}

	subject := "Password Reset Email"
	body := `
		<p>You're receiving this email because you or someone else has requested a password reset for your user account.</p>
		<p>Click the link below to reset your password:
			<a href=%s>Reset Password</a>
		</p>
		<p>Or use this: <br>%s</p>
		<p>If you didn't request this, you can safely ignore this email</p>
		<br/>
		<p>Thanks,</p>
		<p>ProjectWall team</p>
	`

	body = fmt.Sprintf(body, link, link)
	body = mailHead(to, subject) + mailBody(body)

	err := send(body)
	if err != nil {
		err = logger.AppendError(err)
	}

	return err
}
