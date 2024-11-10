package main

import (
	"errors"

	"github.com/muhammadolammi/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		user, err := s.db.GetUserWithName(ctx, s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		if user.Name != s.cfg.CurrentUserName {
			return errors.New("user must be logged in to perform this action")

		}
		return handler(s, c, user)
	}
}
