package main

import (
	"flag"

	"github.com/wpwilson10/caterpillar/internal/news"
	"github.com/wpwilson10/caterpillar/internal/reddit"
	"github.com/wpwilson10/caterpillar/internal/setup"
)

func main() {
	// setup environment configuration
	setup.EnvironmentConfig()

	// get a file for logging
	file := setup.LogFile()
	defer file.Close()

	// setup logger
	setup.Logger(file)

	// run appropriate app
	name, port, app := selectApp()

	if app != nil {
		setup.Application(name)
		setup.RunOnce(port, app)
	}
}

// select app parses input arguments and returns an app's configuration
func selectApp() (string, int, func()) {
	// check command line arguments
	redditBotFlag := flag.Bool("redditBot", false, "RedditBot")
	redditAppFlag := flag.Bool("redditApp", false, "RedditApp")
	newsAppFlag := flag.Bool("newsApp", false, "NewsApp")
	flag.Parse()

	// return appropriate app information
	switch {
	case *redditBotFlag:
		return "RedditBot", setup.EnvToInt("REDDIT_BOT_PORT"), reddit.BotApp
	case *redditAppFlag:
		return "RedditApp", setup.EnvToInt("REDDIT_PORT"), reddit.App
	case *newsAppFlag:
		return "NewsApp", setup.EnvToInt("NEWSPAPER_PORT"), news.App
	}

	// don't do anything on no match
	setup.LogCommon(nil).Fatal("No matching input flag")
	return "", 0, nil
}
