package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

var temp *template.Template

type message struct {
	Key            string
	EncodedMessage string
	DecodedMessage string
	Valid bool
}

type resultStats struct {
	EmptyKey     bool
	EmptyMessage bool
	Message      message
}

func main() {
	fmt.Print("Server should be running now at http://localhost:5978/...")
	temp = template.Must(template.ParseGlob("templates/*.html"))
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", fs))
	http.HandleFunc("/", runIndex)
	http.ListenAndServe(":5978", nil)
}

func runIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp.ExecuteTemplate(w, "index.html", nil)
	} else if r.Method == http.MethodPost {
		rg := regexp.MustCompile(`(\r\n?|\n){2,}`)
		key := rg.ReplaceAllString(strings.ToLower(strings.ReplaceAll(r.FormValue("key"), " ", "")), "")
		encodedMessage := rg.ReplaceAllString(strings.ToLower(r.FormValue("encodedMessage")), " ")
		decodedMessage, valid := decode(key, encodedMessage)
		msg := message{
			Key:            key,
			EncodedMessage: encodedMessage,
			DecodedMessage: decodedMessage,
			Valid: valid,
		}

		temp.ExecuteTemplate(w, "index.html", resultStats{
			len(key) == 0,
			len(strings.ReplaceAll(encodedMessage, " ", "")) == 0,
			msg})
	}
}

func decode(key, text string) (string, bool) {
	alphabet := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	decoderMap := make(map[string]string)
	lastAlphabetIndexVisited := 0
	for _, letter := range key {
		l := string(letter)
		if decoderMap[l] == "" {
			decoderMap[l] = alphabet[lastAlphabetIndexVisited]
			lastAlphabetIndexVisited++
		}

	}

	message := ""
	for _, c := range text {
		letter := string(c)
		if letter == " " {
			message += letter
		} else {
			message += decoderMap[letter]
		}
	}

	return message, len(message) > 0
}
