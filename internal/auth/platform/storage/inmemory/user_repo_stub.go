package inmemory

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/krlspj/go-jwt-auth/internal/auth/domain"
)

var ErrNotUserFound error = errors.New("No user found")

type UserRepositoryStub struct {
	//users []domain.User
	users []userDB
}

func NewUserRepositoryStub() *UserRepositoryStub {
	newUsers := []userDB{
		{id: "1000", username: "John", lastname: "Doe", password: "example"},
		{id: "1001", username: "Mickye", lastname: "Mouse", password: "exmaple"},
	}
	return &UserRepositoryStub{
		users: newUsers,
	}
}

// FindAll return all users stored in memory database
func (s *UserRepositoryStub) FindAll(ctx context.Context) ([]domain.User, error) {
	return toDomainUsers(s.users)
}

func (s *UserRepositoryStub) FindOne(ctx context.Context, id string) (domain.User, error) {
	for _, user := range s.users {
		if user.id == id {
			return user.toDomainUser(), nil
		}
	}
	return domain.User{}, ErrNotUserFound
}
func (s *UserRepositoryStub) CreateUser(ctx context.Context, user domain.User) error {
	udb := toUserDB(user)
	udb.GenerateUUID()
	s.users = append(s.users, udb)

	return nil
}

func (s *UserRepositoryStub) FindOneByField(ctx context.Context, fieldName, fieldValue string) (domain.User, error) {
	fmt.Println("field, value:", fieldName, fieldValue)
	for _, user := range s.users {

		e := reflect.ValueOf(&user).Elem()

		for i := 0; i < e.NumField(); i++ {
			varName := e.Type().Field(i).Name
			//varType := e.Type().Field(i).Type
			varValue := e.Field(i).String()
			//fmt.Printf("%v %v %v\n", varName, varType, varValue)
			if fieldName == varName && fieldValue == varValue {
				return user.toDomainUser(), nil
			}
		}
	}
	return domain.User{}, ErrNotUserFound
}
