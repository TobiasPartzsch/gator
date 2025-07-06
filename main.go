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
	"github.com/tobiaspartzsch/gator/internal/config"
	"github.com/tobiaspartzsch/gator/internal/database"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func main() {
	var err error
	var cfg config.Config
	if cfg, err = config.Read(); err != nil {
		fmt.Fprintf(os.Stderr, "couldn't read config file: %v\n", err)
		os.Exit(1)
	}
	var st state
	st.config = &cfg

	var cmds commands
	cmds.handlers = make(map[string]func(*state, command) error)

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "need at least one argument as the command name")
		os.Exit(1)
	}

	var db *sql.DB
	if db, err = sql.Open("postgres", st.config.DbURL); err != nil {
		fmt.Fprintf(os.Stderr, "error opening database %s: %v\n", st.config.DbURL, err)
		os.Exit(1)
	}
	defer db.Close()
	st.db = database.New(db)

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := command{name: cmdName, args: cmdArgs}
	if err = cmds.run(&st, cmd); err != nil {
		fmt.Fprintf(os.Stderr, "error executing command %s: %v\n", cmdName, err)
		os.Exit(1)
	}

}

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

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("reset should have no arguments")
	}

	err := s.db.DeleteUsers(
		context.Background(),
	)
	if err != nil {
		return fmt.Errorf("error trying to reset database: %w", err)
	}

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

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.name]
	if !exists {
		return fmt.Errorf("no handler registered for command: %s", cmd.name)
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}
