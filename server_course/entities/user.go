package entities

import (
	"context"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (u *User) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)
	if len(u.Email) > 100 {
		problems["too long"] = "email can only be up to and including 100 chars"
	}
	if u.Email == "" {
		problems["no email"] = "messno email set"
	}
	
	return problems
}