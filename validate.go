package main

import (
	"context"
	"errors"
	"regexp"

	fb "github.com/huandu/facebook/v2"
	"google.golang.org/api/idtoken"
)

// Check if code from email valid for use
func isCode(code string) bool {
	if len(code) == 6 {
		return true
	}

	return false
}

// Check email on valid
func IsEmailValid(email string) bool {
	if len(email) < 3 || len(email) > 320 {
		return false
	}

	return regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$").MatchString(email)
}

// Check username on valid
func IsUsernameValid(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return true
	}

	return regexp.MustCompile(`\s`).MatchString(username)
}

// Block of ipa's for Sign In

func isGoogle(id, token string) (string, error) {
	data, err := idtoken.Validate(context.Background(), token, id)

	if err != nil {
		return "", errors.New("0")
	}
	return data.Claims["email"].(string), nil
}

func isFacebook(id, token string) (string, error) {
	data, err := fb.Get("/"+id, fb.Params{
		"fields":       "email",
		"access_token": token,
	})

	if err != nil {
		return "", errors.New("0")
	}
	return data["email"].(string), nil
}
