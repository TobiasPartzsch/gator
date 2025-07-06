package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/tobiaspartzsch/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("login needs exactly one argument")
	}
	_, err := s.db.GetUser(
		context.Background(),
		cmd.args[0],
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Fprintln(os.Stderr, "user does not exist")
			os.Exit(1)
		}
		return fmt.Errorf("error querying user: %w", err)
	}
	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("error trying to set user: %w", err)
	}
	fmt.Printf("user has been set to %s\n", s.config.CurrentUserName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("register needs exactly one argument")
	}

	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.args[0],
		},
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			fmt.Fprintln(os.Stderr, "user already exists")
			os.Exit(1)
		}
		return fmt.Errorf("error trying to create user: %w", err)
	}

	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("could not update config: %w", err)
	}
	fmt.Printf("User '%s' created successfully!\n", user.Name)
	fmt.Printf("Debug: %+v\n", user)
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("users should have no arguments")
	}
	users, err := s.db.GetUsers(
		context.Background(),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Fprintln(os.Stderr, "no users exist")
			os.Exit(1)
		}
		return fmt.Errorf("error querying users: %w", err)
	}

	currentUserName := s.config.CurrentUserName
	for _, user := range users {
		current := ""
		if user.Name == currentUserName {
			current = " (current)"
		}
		fmt.Println("* " + user.Name + current)
	}
	return nil
}
