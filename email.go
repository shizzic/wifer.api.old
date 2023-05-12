package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	mail "github.com/xhit/go-simple-mail/v2"
)

func SendCode(to, code, id string) error {
	server := mail.NewSMTPClient()
	server.Host = config.EMAIL.HOST
	server.Port = config.EMAIL.PORT
	server.Username = config.EMAIL.USERNAME
	server.Password = config.EMAIL.PASSWORD
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	// server.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	smtpClient, err := server.Connect()

	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom("dateshipper <" + config.EMAIL.USERNAME + ">").
		AddTo(to).
		SetSubject("Confirm registration")

	msgUUID, _ := uuid.NewRandom()
	msgID := fmt.Sprintf("<%s@"+config.SELF_DOMAIN_NAME+">", msgUUID.String())
	email.AddHeader("Message-ID", msgID)

	email.SetBody(mail.TextHTML, "<p><h1>Here is a link to sign into dateshipper :)</h1></p><p><a href=\""+config.CLIENT_DOMAIN+"/auth/"+id+"/"+code+"\">Enjoy</a></p>")
	err = email.Send(smtpClient)

	if err != nil {
		return err
	}

	return nil
}

func ContactMe(name, sender, subject, message string) error {
	server := mail.NewSMTPClient()
	server.Host = config.EMAIL.HOST
	server.Port = config.EMAIL.PORT
	server.Username = config.EMAIL.USERNAME
	server.Password = config.EMAIL.PASSWORD
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	// server.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	smtpClient, err := server.Connect()

	if err != nil {
		log.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom("dateshipper <" + config.EMAIL.USERNAME + ">").
		AddTo(config.ADMIN_EMAIL).
		SetSubject(subject)

	msgUUID, _ := uuid.NewRandom()
	msgID := fmt.Sprintf("<%s@"+config.SELF_DOMAIN_NAME+">", msgUUID.String())
	email.AddHeader("Message-ID", msgID)

	email.SetBody(mail.TextHTML, "<p>name: "+name+" - email: "+sender+"</p>"+message)
	err = email.Send(smtpClient)

	if err != nil {
		log.Println(err)
	}

	return nil
}
