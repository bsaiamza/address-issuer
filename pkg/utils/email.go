package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"net/mail"
	"os"
	"strconv"
	"time"

	"address-issuer/pkg/config"
	"address-issuer/pkg/log"

	mailV2 "github.com/xhit/go-simple-mail/v2"
)

func ValidEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("email %s is not valid", email)
	}

	return nil
}

func SendCredentialByEmail(recipientEmail, qrImgName string, qrCode []byte, config *config.Config) error {
	img, _, err := image.Decode(bytes.NewReader(qrCode))
	if err != nil {
		log.ServerError.Printf("Failed to decode qr code: %s", err.Error())
	}

	out, _ := os.Create("./" + qrImgName + ".png")
	defer out.Close()

	err = png.Encode(out, img)
	if err != nil {
		log.ServerError.Printf("Failed to encode png: %s", err.Error())
	}

	htmlBody := `
		<html>
			<head>
				<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
				<title>Physical Address Issuer</title>
			</head>

			<body>
				<p>
					Hi
				</p>
				<p>
					Please scan the following <b>QR Code</b> with your wallet app to get your <b>Physical Address Credential</b>:
					<br/>
					<br/>
					<img src="cid:` + qrImgName + `.png" alt="QR Code" />
				</p>
				<p>
					Regards
					<br/>
					Physical Address Issuer
				</p>
			</body>
		</html>
	`

	server := mailV2.NewSMTPClient()

	//SMTP Server
	server.Host = config.GetSmtpServer()
	smtpPort, err := strconv.Atoi(config.GetSmtpPort())
	if err != nil {
		log.ServerError.Printf("Failed to convert SMTP port: %s", err.Error())
	}
	server.Port = smtpPort
	server.Username = config.GetEmailUsername()
	server.Password = config.GetEmailPassword()
	server.Encryption = mailV2.EncryptionSTARTTLS
	server.Authentication = mailV2.AuthLogin

	//Variable to keep alive connection
	server.KeepAlive = true

	//Timeout for connect to SMTP Server
	server.ConnectTimeout = 10 * time.Second

	//Timeout for send the data and wait respond
	server.SendTimeout = 10 * time.Second

	//SMTP client
	smtpClient, err := server.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect SMTP Client: %s", err.Error())
	}

	//New email simple html with inline and CC
	email := mailV2.NewMSG()

	email.SetFrom(config.GetEmailUsername()).
		AddTo(recipientEmail).
		// AddCc("otherto@example.com").
		SetSubject("Get your Physical Address Credential")

	email.SetBody(mailV2.TextHTML, htmlBody)

	email.Attach(&mailV2.File{FilePath: "./" + qrImgName + ".png", Inline: true})

	// always check error after send
	if email.Error != nil {
		return fmt.Errorf("failed: %s", email.Error)
	}

	// Call Send and pass the client
	err = email.Send(smtpClient)
	if err != nil {
		return fmt.Errorf("failed to send email: %s", err.Error())
	} else {
		log.Info.Println("Email Sent")
	}

	return nil
}

func SendNotificationEmail(recipientEmail string, config *config.Config) error {
	htmlBody := `
		<html>
			<head>
				<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
				<title>Physical Address Issuer</title>
			</head>

			<body>
				<p>
					Hi
				</p>
				<p>
					Please note that your Physical Address Credential has been issued and should be stored in your wallet.
				</p>
				<p>
					Regards
					<br/>
					Physical Address Issuer
				</p>
			</body>
		</html>
	`

	server := mailV2.NewSMTPClient()

	//SMTP Server
	server.Host = config.GetSmtpServer()
	smtpPort, err := strconv.Atoi(config.GetSmtpPort())
	if err != nil {
		log.ServerError.Printf("Failed to convert SMTP port: %s", err.Error())
	}
	server.Port = smtpPort
	server.Username = config.GetEmailUsername()
	server.Password = config.GetEmailPassword()
	server.Encryption = mailV2.EncryptionSTARTTLS
	server.Authentication = mailV2.AuthLogin

	//Variable to keep alive connection
	server.KeepAlive = true

	//Timeout for connect to SMTP Server
	server.ConnectTimeout = 10 * time.Second

	//Timeout for send the data and wait respond
	server.SendTimeout = 10 * time.Second

	//SMTP client
	smtpClient, err := server.Connect()
	if err != nil {
		log.ServerError.Printf("Failed to connect SMTP Client: %s", err.Error())
	}

	//New email simple html with inline and CC
	email := mailV2.NewMSG()

	email.SetFrom(config.GetEmailUsername()).
		AddTo(recipientEmail).
		// AddCc("otherto@example.com").
		SetSubject("Your Physical Address Credential is ready")

	email.SetBody(mailV2.TextHTML, htmlBody)
	// always check error after send
	if email.Error != nil {
		log.ServerError.Println(email.Error)
	}

	// Call Send and pass the client
	err = email.Send(smtpClient)
	if err != nil {
		log.ServerError.Printf("Failed to send email: %s", err.Error())
	} else {
		log.Info.Println("Email Sent")
	}

	return nil
}
