package main

import (
	"crypto/tls"
	"errors"

	h "net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/gomail.v2"
)

// cookie, err := c.Cookie("token")
// c.SetSameSite(4)
// c.SetSameSite(h.SameSiteNoneMode)
// c.SetCookie("token", "here'll be token", 120, "/", "https://wifer-test.ru", true, true)

func SendVerifyEmail(username, to, token, id string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "Wifer <admin@"+domainBack+">")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Confirm registration")
	m.SetBody("text/html", "<p>Hello dear <b>"+username+"</b>. To verify your account, just follow this link<p><a href='"+domainClient+"/ensure/"+username+"/"+token+"/"+id+"'>Verfy your account</a></p></p><p>If you don't signed up on Wifer, than, just follow this link for removing your personal data from service:<br><a href='"+domainClient+"/ensureDelete/"+id+"/"+token+"'>Permanently delete your account</a></p>")
	d := gomail.NewDialer("skvmrelay.netangels.ru", 25, "admin@"+domainBack, emailPass)
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
		{Key: "$set", Value: bson.D{{Key: "active", Value: true}}},
	}); err != nil || r.ModifiedCount == 0 {
		return errors.New("document not updated")
	}

	c.SetSameSite(h.SameSiteNoneMode)
	c.SetCookie("token", EncryptToken(username), 86400*60, "/", domainBack, true, true)
	c.SetCookie("username", username, 86400*60, "/", domainBack, true, true)
	c.SetCookie("id", id, 86400*60, "/", domainBack, true, true)

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

// Send link to new user's email
func SendChangeEmail(oldEmail, newEmail string, c gin.Context) {
	id, _ := c.Cookie("id")
	username, _ := c.Cookie("username")
	token := MakeToken()
	ensure.InsertOne(ctx, bson.D{{Key: "_id", Value: id}, {Key: "token", Value: token}})

	m := gomail.NewMessage()
	m.SetHeader("From", "Wifer <admin@"+domainBack+">")
	m.SetHeader("To", newEmail)
	m.SetHeader("Subject", "Password change")
	m.SetBody("text/html", "<p>Hello dear <b>"+username+"</b>. Just follow link below to confirm your new email.</p><a href='"+domainClient+"/changePassword/"+id+"/"+token+"/"+newEmail+"'>Confirm</a>")
	d := gomail.NewDialer("skvmrelay.netangels.ru", 25, "admin@"+domainBack, emailPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	d.DialAndSend(m)
}

// Send to user link for reset password
func ForgotPassword(to string) error {
	// check if user exists
	var user bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 1})
	if err := users.FindOne(ctx, bson.M{"email": to}, opts).Decode(&user); err != nil {
		return errors.New("0")
	}

	// make new token and ensure that there is no another ensure document in database
	token := MakeToken()
	ensure.DeleteOne(ctx, bson.M{"_id": user["_id"]})

	if _, err := ensure.InsertOne(ctx, bson.D{
		{Key: "_id", Value: user["_id"]},
		{Key: "token", Value: token},
	}); err != nil {
		return errors.New("1")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "Wifer <admin@"+domainBack+">")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "New password")
	m.SetBody("text/html", "<p>Hello there! To change your password, just follow this link<p><a href='"+domainClient+"/newPassword/"+token+"'>Create new password</a></p></p><p>If you didn't try to change your password, than ignore this message :)</a></p>")
	d := gomail.NewDialer("skvmrelay.netangels.ru", 25, "admin@"+domainBack, emailPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return errors.New("2")
	}

	return nil
}
