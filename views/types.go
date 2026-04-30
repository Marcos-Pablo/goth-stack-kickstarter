package views

import (
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/validator"
	"unicode/utf8"
)

func UserInitials(name string) string {
	if name == "" {
		return "?"
	}
	r, _ := utf8.DecodeRuneInString(name)
	return string(r)
}

type User struct {
	Email          string
	Name           string
	ProfilePicture string
}

type SignInForm struct {
	Email        string
	Password     string
	GeneralError string
	validator.Validator
}

func (f *SignInForm) Validate() {
	f.Check(validator.NotBlank(f.Email), "email", "Email is required")
	f.Check(validator.IsEmail(f.Email), "email", "Enter a valid email")
	f.Check(validator.MaxChars(f.Email, 255), "email", "Email is too long")

	f.Check(validator.NotBlank(f.Password), "password", "Password is required")
	f.Check(validator.MinChars(f.Password, 3), "password", "At least 3 characters")
	f.Check(validator.MaxChars(f.Password, 72), "password", "At most 72 characters")
}

type SignUpForm struct {
	Email                string
	Name                 string
	Password             string
	PasswordConfirmation string
	GeneralError         string
	validator.Validator
}

func (f *SignUpForm) Validate() {
	f.Check(validator.NotBlank(f.Email), "email", "Email is required")
	f.Check(validator.IsEmail(f.Email), "email", "Enter a valid email")
	f.Check(validator.MaxChars(f.Email, 255), "email", "Email is too long")

	f.Check(validator.NotBlank(f.Name), "name", "Name is required")
	f.Check(validator.MaxChars(f.Name, 255), "name", "Name is too long")

	f.Check(validator.NotBlank(f.Password), "password", "Password is required")
	f.Check(validator.MinChars(f.Password, 3), "password", "At least 3 characters")
	f.Check(validator.MaxChars(f.Password, 72), "password", "At most 72 characters")
	f.Check(f.Password == f.PasswordConfirmation, "password_confirmation", "Passwords do not match")
}

type ProfileForm struct {
	Email        string
	Name         string
	GeneralError string
	validator.Validator
}

func (f *ProfileForm) Validate() {
	f.Check(validator.NotBlank(f.Email), "email", "Email is required")
	f.Check(validator.IsEmail(f.Email), "email", "Enter a valid email")
	f.Check(validator.MaxChars(f.Email, 255), "email", "Email is too long")

	f.Check(validator.NotBlank(f.Name), "name", "Name is required")
	f.Check(validator.MaxChars(f.Name, 255), "name", "Name is too long")
}

type ChangePasswordForm struct {
	CurrentPassword         string
	NewPassword             string
	NewPasswordConfirmation string
	GeneralError            string
	validator.Validator
}

func (f *ChangePasswordForm) Validate() {
	f.Check(validator.NotBlank(f.CurrentPassword), "current_password", "Current password is required")
	f.Check(validator.MinChars(f.CurrentPassword, 3), "new_password", "At least 3 characters")
	f.Check(validator.MaxChars(f.CurrentPassword, 72), "new_password", "At most 72 characters")

	f.Check(validator.MinChars(f.NewPassword, 3), "new_password", "At least 3 characters")
	f.Check(validator.MaxChars(f.NewPassword, 72), "new_password", "At most 72 characters")

	f.Check(f.NewPassword == f.NewPasswordConfirmation, "new_password_confirmation", "Passwords do not match")
}

type AvatarForm struct {
	GeneralError string
}
