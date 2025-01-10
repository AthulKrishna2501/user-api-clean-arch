package models

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func ValidateSignup(data SignupInput) error {
	if strings.TrimSpace(data.UserName) == "" || strings.TrimSpace(data.Email) == "" || strings.TrimSpace(data.Password) == "" || strings.TrimSpace(data.PhoneNumber) == "" {
		return errors.New(ErrRequiredFieldsEmpty)
	}

	if err := ValidateEmail(data.Email); err != nil {
		return err
	}

	if err := ValidatePhoneNumber(data.PhoneNumber); err != nil {
		return err
	}

	if err := ValidatePassword(data.Password); err != nil {
		return err
	}

	return nil
}

func ValidateEmail(email string) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		return errors.New("error while validating email format")
	}
	if !matched {
		return errors.New(ErrInvalidEmailFormat)
	}
	return nil
}

func ValidatePhoneNumber(phoneNumber string) error {
	phoneRegex := `^\+?[1-9]\d{1,14}$`
	matched, err := regexp.MatchString(phoneRegex, phoneNumber)
	if err != nil {
		return errors.New("error while validating phone number")
	}
	if !matched {
		return errors.New(ErrInvalidPhoneNumber)
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength || len(password) > MaxPasswordLength {
		return fmt.Errorf(ErrPasswordLength, MinPasswordLength, MaxPasswordLength)
	}

	complexityRegex := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).+$`
	matched, err := regexp.MatchString(complexityRegex, password)
	if err != nil {
		return errors.New("error while validating password complexity")
	}

	if !matched {
		return errors.New(ErrPasswordComplexity)
	}

	return nil
}
