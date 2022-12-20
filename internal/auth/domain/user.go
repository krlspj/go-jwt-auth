package domain

import "context"

type User struct {
	id          string
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

type UserRepo interface {
	FindAll(ctx context.Context) ([]User, error)
}
