package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	internal "github.com/muhammadolammi/gator/internal/config"
	"github.com/muhammadolammi/gator/internal/database"
)

var ctx = context.Background()

func main() {
	// fmt.Println(internal.Read())
	cfg := internal.Read()

	// open a db connectuion, with the dburl in config
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Panic(err)

	}
	dbQueries := database.New(db)
	st := state{
		cfg: &cfg,
		db:  dbQueries,
	}
	commands := commands{
		cmds: map[string]func(*state, command) error{},
	}
	// register commands
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", middlewareLoggedIn(handlerGetFeeds))
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", middlewareLoggedIn(handlerFollowing))
	commands.register("unfollow", middlewareLoggedIn(handlerUnFollow))
	commands.register("browse", middlewareLoggedIn(handlerBrowse))

	args := os.Args
	if len(args) < 2 {
		log.Println("command and args expected")
		return

	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	err = commands.run(&st, cmd)
	if err != nil {
		log.Println(err)
		return

	}

}
