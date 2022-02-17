package main

import (
	"crypto/tls"
	"errors"

	h "net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/gomail.v2"
)

var domain string = "https://wifer.ru"

// cookie, err := c.Cookie("token")
// c.SetSameSite(4)
// c.SetSameSite(h.SameSiteNoneMode)
// c.SetCookie("token", "here'll be token", 120, "/", "https://wifer-test.ru", true, true)

func SendVerifyEmail(username, to, token, id string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "Wifer <admin@wifer-test.ru>")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Confirm registration")
	m.SetBody("text/html", "<p>Hello dear <b>"+username+"</b>. To verify your account, just follow this link<p><a href='"+domain+"/ensure/"+username+"/"+token+"/"+id+"'>Verfy your account</a></p></p><p>If you don't signed up on Wifer, than, just follow this link for removing your personal data from service:<br><a href='"+domain+"/ensureDelete/"+id+"/"+token+"'>Permanently delete your account</a></p>")
	d := gomail.NewDialer("skvmrelay.netangels.ru", 25, "admin@wifer-test.ru", "jukdNRaVWf3Fvmg")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// Activate user's account
func EnsureEmail(username, token, id string, c gin.Context) error {
	if r, err := ensure.DeleteOne(ctx, bson.M{"_id": id, "token": token}); err != nil || r.DeletedCount == 0 {
		return errors.New("not found")
	}

	if r, err := users.UpdateOne(ctx, bson.M{"_id": id}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "status", Value: true}}},
	}); err != nil || r.ModifiedCount == 0 {
		return errors.New("document not updated")
	}

	c.SetSameSite(h.SameSiteNoneMode)
	c.SetCookie("token", EncryptToken(username), 120, "/", "wifer-test.ru", true, true)
	c.SetCookie("username", username, 120, "/", "wifer-test.ru", true, true)
	c.SetCookie("id", id, 120, "/", "wifer-test.ru", true, true)

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

// Tell user that his account has been deleted
func InfoAboutDelete(email, username string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "Wifer <admin@wifer-test.ru>")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Account removing")
	m.SetBody("text/html", "<p>Hello dear <b>"+username+"</b>. Your account deleted successfully!")
	d := gomail.NewDialer("skvmrelay.netangels.ru", 25, "admin@wifer-test.ru", "jukdNRaVWf3Fvmg")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	d.DialAndSend(m)
}

// Tell user that his password was changed
func InfoAboutPasswordChange(email, username string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "Wifer <admin@wifer-test.ru>")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Password change")
	m.SetBody("text/html", "<p>Hello dear <b>"+username+"</b>. Your password was successfully changed")
	d := gomail.NewDialer("skvmrelay.netangels.ru", 25, "admin@wifer-test.ru", "jukdNRaVWf3Fvmg")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	d.DialAndSend(m)
}

// Send link to new user's email
func SendChangeEmail(oldEmail, newEmail string, c gin.Context) {
	id, _ := c.Cookie("id")
	username, _ := c.Cookie("username")
	token := MakeToken()
	ensure.InsertOne(ctx, bson.D{{Key: "_id", Value: id}, {Key: "token", Value: token}})

	m := gomail.NewMessage()
	m.SetHeader("From", "Wifer <admin@wifer-test.ru>")
	m.SetHeader("To", newEmail)
	m.SetHeader("Subject", "Password change")
	m.SetBody("text/html", "<p>Hello dear <b>"+username+"</b>. Just follow link below to confirm your new email.</p><a href='"+domain+"/changePassword/"+id+"/"+token+"/"+newEmail+"'>Confirm</a>")
	d := gomail.NewDialer("skvmrelay.netangels.ru", 25, "admin@wifer-test.ru", "jukdNRaVWf3Fvmg")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	d.DialAndSend(m)
}
