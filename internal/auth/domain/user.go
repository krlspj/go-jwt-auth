package domain

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id string
	//Xname       string `validate:"required"` -> To validate this need to be public (X) not private (x)
	name        string
	lastname    string
	password    string
	token       string
	createdTime string
	updatedTime string
}

// Getters and Setters

func (u *User) ID() string {
	return u.id
}
func (u *User) SetID(id string) {
	u.id = id
}

func (u *User) Name() string {
	return u.name
}
func (u *User) SetName(name string) {
	u.name = name
}

func (u *User) Lastname() string {
	return u.lastname
}
func (u *User) SetLastname(lastname string) {
	u.lastname = lastname
}

func (u *User) Password() string {
	return u.password
}
func (u *User) SetPassword(password string) {
	u.password = password
}

func (u *User) Token() string {
	return u.token
}
func (u *User) SetToken(token string) {
	u.token = token
}

// Other methods

// HasPassword hashes the password from the user
func (u *User) HashPassword() error {
	hashedPass, err := hashPassword(u.Password())
	if err != nil {
		return err
	}
	u.SetPassword(hashedPass)
	return nil
}

// VerifyPassword, returns an error if the provided password is not the one for this user
func (u *User) VerifyPassword(providedPass string) error {
	if err := verifyPassword(u.Password(), providedPass); err != nil {
		return err
	}
	return nil
}

//func (u *User) ValidateUser() error {
//	fmt.Println(" validateUser")
//	validate := validator.New()
//	if err := validate.Struct(u); err != nil {
//		return err
//	}
//	return nil
//}

// HashPassword returns an encrypted passwor
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // this hash generation already uses a salt in it
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// VerifyPassword returns true if passwords maches. When not, false and a message
func verifyPassword(userPassword string, providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	if err != nil {
		return err
	}
	return nil
}

type UserRepo interface {
	FindAll(ctx context.Context) ([]User, error)
	// FindOne use the mongo id to find the record
	FindOne(ctx context.Context, id string) (User, error)
	FindOneByField(ctx context.Context, fieldName, fieldValue string) (User, error)
	CreateUser(ctx context.Context, user User) error
}
