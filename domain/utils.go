package domain

import "regexp"

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidName(name string) bool {
	nameRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
	return nameRegex.MatchString(name)
}

func IsValidPhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^-?\d+$`)
	return phoneRegex.MatchString(phone)
}
