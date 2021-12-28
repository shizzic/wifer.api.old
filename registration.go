package main

import (
	"errors"
	"math/rand"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Struct of new user
type registrat struct {
	Username string `form:"username"`
	Code     string `form:"countryCode"`
	Email    string `form:"email"`
	Password string `form:"password"`
	Title    string `form:"title"`
	Age      uint8  `form:"age"`
	Height   uint8  `form:"height"`
	Weight   uint8  `form:"weight"`
	Body     uint8  `form:"body"`
}

// Check each value from POST form on correct format.
// Return either success message or error
func Registration(data registrat) error {
	if isUsernameValid(data.Username) {
		return errors.New("incorrect username")
	}

	if !isEmailValid(data.Email) {
		return errors.New("incorrect email")
	}

	if !isPasswordValid(&data.Password) {
		return errors.New("incorrect password")
	}

	if isTitleValid(data.Title) {
		return errors.New("incorrect title")
	}

	if isAgeValid(data.Age) {
		return errors.New("incorrect age")
	}

	if isHeightValid(data.Height) {
		return errors.New("incorrect height")
	}

	if isWeightValid(data.Weight) {
		return errors.New("incorrect weight")
	}

	if isBodyValid(data.Body) {
		return errors.New("incorrect body")
	}

	if _, err := users.InsertOne(ctx, bson.D{
		{Key: "username", Value: data.Username},
		{Key: "email", Value: data.Email},
		{Key: "password_hash", Value: data.Password},
		{Key: "title", Value: data.Title},
		{Key: "age", Value: data.Age},
		{Key: "body", Value: data.Body},
		{Key: "height", Value: data.Height},
		{Key: "weight", Value: data.Weight},
		{Key: "premium", Value: false},
		{Key: "status", Value: false},
	}); err != nil {
		return errors.New("document not inserted")
	}

	token := generateTokenForEmail()

	if _, err := ensure.InsertOne(ctx, bson.D{
		{Key: "_id", Value: data.Username},
		{Key: "token", Value: token},
	}); err != nil {
		return errors.New("document not inserted")
	}

	if err := SendVerifyEmail(data.Username, data.Email, token); err != nil {
		return errors.New("couldn't send message to your email")
	}

	return nil
}

// Check email on valid
func isEmailValid(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}

	return regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$").MatchString(email)
}

// Check password on valid length and do hash
func isPasswordValid(password *string) bool {
	if len(*password) < 8 || len(*password) > 128 {
		return false
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(*password), 8)
	*password = string(hashed)
	return true
}

// Check username on valid
func isUsernameValid(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return true
	}

	return regexp.MustCompile(`\s`).MatchString(username)
}

// Check title on valid
func isTitleValid(title string) bool {
	return len(title) > 150
}

// Check age on valid
func isAgeValid(age uint8) bool {
	return age < 18 || age > 100
}

// Check height on valid
func isHeightValid(height uint8) bool {
	return height < 140 || height > 220
}

// Check weight on valid
func isWeightValid(weight uint8) bool {
	return weight < 30 || weight > 220
}

// Check body on valid
func isBodyValid(body uint8) bool {
	return body < 1 || body > 7
}

// make token for EnsureEmail
func generateTokenForEmail() string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 64)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}

	return string(b)
}
