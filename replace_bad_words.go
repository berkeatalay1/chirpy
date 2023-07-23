package main

import "strings"

func replace_bad_words(str string) string {
	var wrongWords [3]string
	wrongWords[0] = "kerfuffle"
	wrongWords[1] = "sharbert"
	wrongWords[2] = "fornax"
	words := strings.Split(str, " ")
	for i, word := range words {
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "fornax" {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
