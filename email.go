package main

import (
	"crypto/tls"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/gomail.v2"
)

func SendCode(to, code, id string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "Wifer <admin@"+domainBack+">")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Confirm registration")
	m.SetBody("text/html", "<p><h1>Here is a link to sign in into wifer :)</h1></p><p><a href=\"http://192.168.1.106:8080/auth/"+id+"/"+code+"\">Enjoy</a></p>")
	d := gomail.NewDialer("skvmrelay.netangels.ru", 25, "admin@"+domainBack, emailPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// Delete all user's data, if it triggered by link from email
func EnsureDelete(id, token string) error {
	if r, err := ensure.DeleteOne(ctx, bson.M{"_id": id, "token": token}); err != nil || r.DeletedCount == 0 {
		return errors.New("ensure hasn't deleted")
	}

	if r, err := users.DeleteOne(ctx, bson.M{"_id": id}); err != nil || r.DeletedCount == 0 {
		return errors.New("user hasn't deleted")
	}

	return nil
}
