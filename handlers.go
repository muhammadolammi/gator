package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
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
	if len(cmd.args) != 1 {
		return errors.New("agg command expect a time_between_reqs as the args")
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("error parsing time_between_reqs arg to a time.Duration.  err: %v", err)
	}
	fmt.Printf("Collecting feeds every %v", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return errors.New("addfeed command expect a feed name and feed url as the args")
	}
	user, err := s.db.GetUserWithName(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}
	// follow the feed for the user automatically
	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    feed.UserID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println(feed)

	return nil
}

func handlerGetFeeds(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("%s %s %s ", feed.Name, feed.Url, feed.Name)

		user, err := s.db.GetUser(ctx, feed.UserID)
		if err != nil {
			return err
		}
		fmt.Println(user.Name)
	}
	return nil

}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("follow command expect a feed url as the args")
	}
	// user, err := s.db.GetUserWithName(ctx, s.cfg.CurrentUserName)
	// if err != nil {
	// 	return err
	// }
	feed, err := s.db.GetFeedWithUrl(ctx, cmd.args[0])
	if err != nil {
		return err
	}
	feed_follow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println(feed_follow.FeedName)
	fmt.Println(feed_follow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	user_followings, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}
	for _, user_following := range user_followings {
		fmt.Println(user_following.FeedName)
	}
	return nil
}

func handlerUnFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("unfollow command expect a feed_url as the args")
	}
	err := s.db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		Username: user.Name,
		FeedUrl:  cmd.args[0],
	})

	return err
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32
	limit = 2

	if len(cmd.args) == 1 {

		l, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			limit = 2
		}
		limit = int32(l)
	}

	posts, err := s.db.GetPostsForUser(ctx, database.GetPostsForUserParams{
		ID:    user.ID,
		Limit: limit,
	})

	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Println(post)
	}

	return nil
}

func handlerHelp(s *state, cmd command) error {
	fmt.Print(`These are all available functions
	1. gator users
		This will list all users indicating the logged in user.
	2. gator register <user>
		This will register and auto login a  new user
	3. gator login <user>
		This will login to the provided user
	4. gator agg <time_between_reqs >
		run the aggregator that fetch feed every time_between_reqs 
	5. gator reset 
	This will reset and delete all data on the aggregator
	6. gator feeds
		This list all available feeds on the aggregator
	7. gator addfeed <feed name> <feed url>
		This add a new feed to the aggregator feeds
	8. gator follow <feed url> 
		This command  makes the logged in user follow the feed
	9. gator unfollow <feed url> 
		This command  makes the logged in user unfollow the feed
	10. gator following
		This list all feeds the user is following
	11. gator browse
		This list all posts for the current user

	12. gator help
		list all commands and functionality
`)
	fmt.Println()
	return nil
}
