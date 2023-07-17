package domain

import (
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ValidPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
