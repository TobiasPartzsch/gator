package main

import (
	"fmt"

	"github.com/tobiaspartzsch/gator/internal/database"
)

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
	if feed.LastFetchedAt.Valid {
		fmt.Printf("* LastFetchedAt: %v\n", feed.LastFetchedAt.Time)
	} else {
		fmt.Printf("* LastFetchedAt: never\n")
	}
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}
