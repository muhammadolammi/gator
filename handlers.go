package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/muhammadolammi/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("login command expect a username as the args")
	}

	userExist, err := s.db.UserExists(ctx, cmd.args[0])
	if err != nil {
		return err
	}
	if !userExist {
		os.Exit(1)
	}
	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("user set successfully")
	return nil
}

func handlerRegister(s *state, cmd command) error {

	if len(cmd.args) != 1 {
		return errors.New("register command expect a name as the args")
	}
	userExist, err := s.db.UserExists(ctx, cmd.args[0])
	if err != nil {
		return err
	}
	if userExist {
		os.Exit(1)
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{

		Name:      cmd.args[0],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Println("user created")
	fmt.Println(user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(ctx)
	return err
}
func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("%v (current)\n", user.Name)
		} else {
			fmt.Println(user.Name)
		}

	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	rss, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(rss)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return errors.New("addfeed command expect a feed name and feed url as the args")
	}
	cUser, err := s.db.GetUserWithName(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    cUser.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println(feed)

	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf(feed.Name)
		fmt.Printf(feed.Url)
		fmt.Printf(feed.Name)
		user, err := s.db.GetUser(ctx, feed.UserID)
		if err != nil {
			return err
		}
		fmt.Println(user.Name)
	}
	return nil

}
