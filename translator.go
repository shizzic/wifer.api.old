package main

import gt "github.com/bas24/googletranslatefree"

func Translate(text, lang string) (string, error) {
	result, err := gt.Translate(text, "auto", lang)

	if err != nil {
		return "", err
	} else {
		return result, nil
	}
}
