package service

import "context"

type User struct {
	userRepo UserRepo
}

func NewUser(userRepo UserRepo) *User {
	return &User{
		userRepo: userRepo,
	}
}

func (u *User) Create(ctx context.Context) (int64, error) {
	return u.userRepo.CreateUser(ctx)
}
