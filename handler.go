package main

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/palantir/go-githubapp/githubapp"
)

type Handler struct {
	githubapp.ClientCreator
}

func (h *Handler) Handles() []string {
	return []string{"issue_comment"}
}

func (h *Handler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.IssueCommentEvent

	installationID := githubapp.GetInstallationIDFromEvent(&event)
	client, err := h.NewInstallationClient(installationID)
	if err != nil {
		return err
	}
	_ = client

	return nil
}
