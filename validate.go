package main

import (
	"context"
	"errors"
	"regexp"

	fb "github.com/huandu/facebook/v2"
	"google.golang.org/api/idtoken"
)

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

// --------------------------------------------------------------------

// Check if code from email valid for use
func isCode(code string) bool {
	return len(code) == 6
}

func IsEmailValid(value string) bool {
	if len(value) < 3 || len(value) > 320 {
		return false
	}

	return regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$").MatchString(value)
}

func IsUsernameValid(value string) bool {
	len := len(value)

	if len < 1 || len > 20 {
		return false
	}

	res := regexp.MustCompile(`\s`).MatchString(value)

	return !res
}

func IsTitleValid(value string) bool {
	return len(value) <= 150
}

func IsAboutValid(value string) bool {
	return len(value) <= 1500
}

func IsSmokeValid(value int) bool {
	if value < 0 || value > 3 {
		return false
	}

	return true
}

func IsDrinkValid(value int) bool {
	if value < 0 || value > 3 {
		return false
	}

	return true
}

func IsBodyValid(value int) bool {
	if value < 0 || value > 7 {
		return false
	}

	return true
}

func IsEthnicityValid(value int) bool {
	if value < 0 || value > 9 {
		return false
	}

	return true
}

func IsIncomeValid(value int) bool {
	if value < 0 || value > 12 {
		return false
	}

	return true
}

func IsIndustryValid(value int) bool {
	if value < 0 || value > 17 {
		return false
	}

	return true
}

func IsPreferValid(value int) bool {
	if value < 0 || value > 3 {
		return false
	}

	return true
}

func IsSexValid(value int) bool {
	if value < 0 || value > 2 {
		return false
	}

	return true
}

func IsAgeValid(value int) bool {
	if value < 0 || value > 80 {
		return false
	}

	return true
}

func IsHeightValid(value int) bool {
	if value < 0 || value > 220 {
		return false
	}

	return true
}

func IsWeightValid(value int) bool {
	if value < 0 || value > 220 {
		return false
	}

	return true
}

func IsChildrenValid(value int) bool {
	if value < 0 || value > 9 {
		return false
	}

	return true
}
