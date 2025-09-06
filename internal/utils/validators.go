package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func ValidateMobileNumber(mobile string) bool {
	mobile = strings.TrimSpace(mobile)

	pattern := `^[6-9][0-9]{9}$`

	re := regexp.MustCompile(pattern)
	return re.MatchString(mobile)
}

func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return false
	}
	return true
}

func ValidatePassword(password string) bool {
	var hasLower, hasDigit, hasSpecial bool

	if len(password) < 12 {
		return false
	}

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasLower || !hasDigit || !hasSpecial {
		return false
	}

	return true
}
