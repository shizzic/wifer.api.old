package main

import (
	"crypto/tls"
	"errors"

	h "net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/gomail.v2"
)

// cookie, err := c.Cookie("token")
// c.SetSameSite(4)
// c.SetSameSite(h.SameSiteNoneMode)
// c.SetCookie("token", "here'll be token", 120, "/", "https://wifer-test.ru", true, true)

func SendVerifyEmail(username, to, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "admin@wifer-test.ru")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Confirm registration")
	m.SetBody("text/html", "<p>Hello dear <b>"+username+"</b>. To verify your account, just follow this link<p><a href='https://wifer.ru/ensure/"+username+"/"+token+"'>Verfy your account</a></p></p><p>If you don't signed up on Wifer, than, just follow this link for removing your personal data from service:<br><a href='https://wifer.ru/ensureDelete/"+username+"/"+token+"'>Permanently delete your data</a></p>")
	d := gomail.NewDialer("skvmrelay.netangels.ru", 25, "admin@wifer-test.ru", "jukdNRaVWf3Fvmg")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// Activate user's account
func EnsureEmail(username, token string, c gin.Context) error {
	if r, err := ensure.DeleteOne(ctx, bson.M{"_id": username, "token": token}); err != nil || r.DeletedCount == 0 {
		return errors.New("not found")
	}

	if r, err := users.UpdateOne(ctx, bson.M{"username": username}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "status", Value: true}}},
	}); err != nil || r.ModifiedCount == 0 {
		return errors.New("document not updated")
	}

	c.SetSameSite(h.SameSiteNoneMode)
	c.SetCookie("token", EncryptToken(username), 120, "/", "wifer-test.ru", true, true)
	c.SetCookie("username", username, 120, "/", "wifer-test.ru", true, true)

	return nil
}

// Delete all user's data, if it triggered by link from email
func EnsureDelete(username, token string) error {
	if r, err := ensure.DeleteOne(ctx, bson.M{"_id": username, "token": token}); err != nil || r.DeletedCount == 0 {
		return errors.New("error")
	}

	if _, err := about.DeleteOne(ctx, bson.M{"_id": username}); err != nil {
		return errors.New("error")
	}

	if r, err := users.DeleteOne(ctx, bson.M{"username": username}); err != nil || r.DeletedCount == 0 {
		return errors.New("error")
	}

	return nil
}

// Tell user that his account has been deleted
func InfoAboutDelete(email, username string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "admin@wifer-test.ru")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Account removing")
	m.SetBody("text/html", "<p>Hello dear <b>"+username+"</b>. Your account deleted successfully!")
	d := gomail.NewDialer("skvmrelay.netangels.ru", 25, "admin@wifer-test.ru", "jukdNRaVWf3Fvmg")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	d.DialAndSend(m)
}
