package main

import (
	"context"
	"fmt"
)

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
