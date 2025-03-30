package utils

import (
	"strings"
	"regexp"
	"html"
	"golang.org/x/crypto/bcrypt"
)

func HashPass(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePasswords(hashed string, plain []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), plain)
	return err == nil
}

func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	input = html.EscapeString(input)	// Changes the HTML characters such as <, >, etc so it cant run a script
	re := regexp.MustCompile(`[^\w\s@.-]`)	// Removes unallowed characters such as %, $, &, etc
	return re.ReplaceAllString(input, "")	// Usees the regex from before "re"
}
