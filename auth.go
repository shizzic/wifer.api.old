package main

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// 33 to 126. Space not allowed
func EncryptToken(username string) (token string) {
	for i, char := range username {
		if char%2 == 0 {
			token += string(char - 1)
		} else {
			token += string(char + 1)
		}

		rand.Seed(time.Now().UnixNano())

		b := make([]rune, i)
		for k := range b {
			b[k] = letters[rand.Intn(len(letters))]
		}

		token += string(b)
	}

	return
}

// 30 ms speed average with simple request
func DecryptToken() (username string) {
	token := "lpjsRidNrjjHKJSdTywosgbNuRld"
	key := 0
	minus := 0

	for i, char := range token {
		if key == i {
			if char%2 == 0 {
				username += string(char - 1)
			} else {
				username += string(char + 1)
			}

			key += minus + 1
			minus += 1
		}
	}

	return
}
