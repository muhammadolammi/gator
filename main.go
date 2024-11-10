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
	commands.register("addfeed", handlerAddFeed)
	commands.register("feeds", handlerGetFeeds)
	commands.register("follow", handlerFollow)
	commands.register("following", handlerFollowing)

	args := os.Args
	if len(args) < 2 {
		log.Panic("command and args expected")

	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	err = commands.run(&st, cmd)
	if err != nil {
		log.Panic(err)

	}

}
