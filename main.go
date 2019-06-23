package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/ctrlaltdel121/configor"
	"github.com/fanatic/git-training/handlers"
	"github.com/gregjones/httpcache"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Github githubapp.Config `yaml:"github"`
}

func main() {
	var cfg Config
	if err := configor.Load(&cfg, "config.yml"); err != nil {
		logrus.Fatalf("Error loading config: %s\n", err)
	}

	cfg.Github.App.IntegrationID, _ = strconv.Atoi(os.Getenv("INTEGRATION_ID"))
	cfg.Github.App.WebhookSecret = os.Getenv("WEBHOOK_SECRET")
	cfg.Github.OAuth.ClientID = os.Getenv("GITHUB_CLIENT_ID")
	cfg.Github.OAuth.ClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	cfg.Github.App.PrivateKey = os.Getenv("GITHUB_PRIVATE_KEY")

	cc, err := githubapp.NewDefaultCachingClientCreator(
		cfg.Github,
		githubapp.WithClientUserAgent("git-training/0.0.1"),
		githubapp.WithClientCaching(false, func() httpcache.Cache { return httpcache.NewMemoryCache() }),
	)
	if err != nil {
		logrus.Fatalf("Error creating client creator: %s\n", err)
	}

	webhookHandler := githubapp.NewDefaultEventDispatcher(
		cfg.Github,
		&handlers.IssuesHandler{ClientCreator: cc},
		&handlers.CreateHandler{ClientCreator: cc},
	)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	loggingHandler := hlog.NewHandler(logger)(webhookHandler)

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), loggingHandler); err != nil {
		logrus.Fatalf("Error creating client creator: %s\n", err)
	}
}
