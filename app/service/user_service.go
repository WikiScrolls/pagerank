package service

import (
	"context"

	g "github.com/gorse-io/gorse-go"
)

type UserService struct {
	gorse *g.GorseClient
}

func NewUserService(
	gorse *g.GorseClient,
) *UserService {
	return &UserService{gorse: gorse}
}

func (s *UserService) RegisterUser(ctx context.Context, userId string, interests []string) error {
	_, err := s.gorse.InsertUser(ctx, g.User{
		UserId: userId,
		Labels: interests,
	})

	return err
}
