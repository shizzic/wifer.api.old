package main

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

// cookie, err := c.Cookie("token")
// c.SetSameSite(4)
// c.SetCookie("token", "here'll be token", 120, "/", "https://wifer-test.ru", true, true)

const mailHost = "skvmrelay.netangels.ru"
const mailPort = 25
const mailFrom = "admin@wifer-test.ru"
const mailUser = "admin@wifer-test.ru"
const mailPass = "jukdNRaVWf3Fvmg"

func SendVerifyEmail(username, to string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mailFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Confirm registration")
	m.SetBody("text/html", "<p>Hello dear <b>"+username+"</b>.</p><p>To verify your email, just follow this link:<br>"+
		"<a href='http://wifer.ru/token'>Verfy</a></p>")
	d := gomail.NewDialer(mailHost, mailPort, mailUser, mailPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
