package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/tobiaspartzsch/gator/internal/config"
	"github.com/tobiaspartzsch/gator/internal/database"
)

type state struct {
	config *config.Config
	db     *database.Queries
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
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddfeed))
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerListFeedFollows))

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
