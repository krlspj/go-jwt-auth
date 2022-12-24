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
	isAdmin     string
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
func (u *User) CreatedTime() string {
	return u.createdTime
}
func (u *User) SetCreatedTime(ct string) {
	u.createdTime = ct
}
func (u *User) UpdatedTime() string {
	return u.updatedTime
}
func (u *User) SetUpdatedTime(ut string) {
	u.updatedTime = ut
}

func (u *User) IsAdmin() string {
	return u.isAdmin
}
func (u *User) SetIsAdmin(a string) {
	u.isAdmin = a
}

// Other methods

// HasPassword hashes and stores the password in user.password
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

// hashPassword returns an encrypted passwor
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // this hash generation already uses a salt in it
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// verifyPassword returns an error if passwords not match
func verifyPassword(hashPass string, providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

type AuthzUserRepo interface {
	FindOne(ctx context.Context, id string) (User, error)
	//FindAll(ctx context.Context) ([]User, error)
	// FindOne use the mongo id to find the record
	//FindOneByField(ctx context.Context, fieldName, fieldValue string) (User, error)
	//InsertUser(ctx context.Context, user User) error
}
