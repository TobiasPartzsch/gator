package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/tobiaspartzsch/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				fmt.Fprintf(os.Stderr, "user %s does not exist\n", s.config.CurrentUserName)
				os.Exit(1)
			}
			return fmt.Errorf("error querying user: %w", err)
		}

		return handler(s, cmd, user)
	}

}
