package entities

import (
	"context"
	"encoding/json"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID               int    `json:"id"`
	Email            string `json:"email"`
	Password         string `json:"password,omitempty"`
	ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
	Token            string `json:"token"`
}

func (u *User) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)
	if len(u.Email) > 100 {
		problems["too long"] = "email can only be up to and including 100 chars"
	}
	if u.Email == "" {
		problems["no email"] = "no email set"
	}

	return problems
}

func (u *User) ValidPassword(expectedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(expectedPassword), []byte(u.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (u *User) MarshalJSONCustom() ([]byte, error) {
	copyUser := u
	copyUser.Password = ""
	copyUser.ExpiresInSeconds = 0
	return json.Marshal(copyUser)
}

func (u *User) EncryptPassword() (User, error) {
	copyUser := u
	bytePassword, err := bcrypt.GenerateFromPassword([]byte(copyUser.Password), 5)
	copyUser.Password = string(bytePassword)
	return *copyUser, err
}
