package handlers

import (
	"context"
	"encoding/json"

	"github.com/sirupsen/logrus"

	"github.com/google/go-github/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
)

type IssuesHandler struct {
	githubapp.ClientCreator
}

func (h *IssuesHandler) Handles() []string {
	return []string{"issues"}
}

func (h *IssuesHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.IssuesEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse issue event payload")
	}

	switch event.GetAction() {
	case "opened":
		logrus.Info("Handling opened")
	default:
		logrus.Infof("Handling %s", event.GetAction())
	}

	installationID := githubapp.GetInstallationIDFromEvent(&event)
	client, err := h.NewInstallationClient(installationID)
	if err != nil {
		return err
	}
	_ = client

	return nil
}
