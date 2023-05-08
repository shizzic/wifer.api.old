package main

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/google/uuid"
	mail "github.com/xhit/go-simple-mail/v2"
)

func SendCode(to, code, id string) error {
	server := mail.NewSMTPClient()
	server.Host = "skvmrelay.netangels.ru"
	server.Port = 25
	server.Username = "admin@" + SELF_DOMAIN_NAME
	server.Password = EMAIL_PASSWORD
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	smtpClient, err := server.Connect()

	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom("luckriza <admin@" + SELF_DOMAIN_NAME + ">").
		AddTo(to).
		SetSubject("Confirm registration")

	msgUUID, _ := uuid.NewRandom()
	msgID := fmt.Sprintf("<%s@"+SELF_DOMAIN_NAME+">", msgUUID.String())
	email.AddHeader("Message-ID", msgID)

	fmt.Println(msgID)

	email.SetBody(mail.TextHTML, "<p><h1>Here is a link to sign in into luckriza :)</h1></p><p><a href=\""+CLIENT_DOMAIN+"/auth/"+id+"/"+code+"\">Enjoy</a></p>")
	err = email.Send(smtpClient)

	if err != nil {
		return err
	}

	return nil
}

func ContactMe(name, sender, subject, message string) error {
	server := mail.NewSMTPClient()
	server.Host = "skvmrelay.netangels.ru"
	server.Port = 25
	server.Username = "admin@" + SELF_DOMAIN_NAME
	server.Password = EMAIL_PASSWORD
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	smtpClient, err := server.Connect()

	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom("luckriza <admin@" + SELF_DOMAIN_NAME + ">").
		AddTo("kotcich@gmail.com").
		SetSubject(subject)

	msgUUID, _ := uuid.NewRandom()
	msgID := fmt.Sprintf("<%s@"+SELF_DOMAIN_NAME+">", msgUUID.String())
	email.AddHeader("Message-ID", msgID)

	fmt.Println(msgID)

	email.SetBody(mail.TextHTML, "<p>name: "+name+" - email: "+sender+"</p>"+message)
	err = email.Send(smtpClient)

	if err != nil {
		return err
	}

	return nil
}
