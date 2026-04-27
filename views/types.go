package views

import "github.com/Marcos-Pablo/goth-stack-kickstarter/internal/validator"
type User struct {
	Email string
	Name  string
}

type SignInForm struct {
	Email        string
	Password     string
	GeneralError string
	validator.Validator
}

type SignUpForm struct {
	Email                string
	Name                 string
	Password             string
	PasswordConfirmation string
	GeneralError         string
	validator.Validator
}
