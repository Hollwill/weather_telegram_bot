package models

import "context"

type UserRepository interface {
	Save(ctx context.Context, u *User) error
	FindByID(ctx context.Context, id int) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
}
