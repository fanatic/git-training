package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/ctrlaltdel121/configor"
	"github.com/gregjones/httpcache"
	"github.com/palantir/go-baseapp/baseapp"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Server baseapp.HTTPConfig `yaml:"server"`
	Github githubapp.Config   `yaml:"github"`

	AppConfig struct {
		PullRequestPreamble string `yaml:"pull_request_preamble"`
	} `yaml:"app_configuration"`
}

func main() {
	var cfg Config
	if err := configor.Load(&cfg, "config.yml"); err != nil {
		logrus.Fatalf("Error loading config: %s\n", err)
	}

	cfg.Github.App.IntegrationID, _ = strconv.Atoi(os.Getenv("INTEGRATION_ID"))
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

	handler := &Handler{
		ClientCreator: cc,
	}

	webhookHandler := githubapp.NewDefaultEventDispatcher(cfg.Github, handler)

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), webhookHandler); err != nil {
		logrus.Fatalf("Error creating client creator: %s\n", err)
	}
}
