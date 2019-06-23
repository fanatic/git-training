package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/google/go-github/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
)

type PullRequestHandler struct {
	githubapp.ClientCreator
}

func (h *PullRequestHandler) Handles() []string {
	return []string{"pull_request"}
}

func (h *PullRequestHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.PullRequestEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse pull_request event payload")
	}

	logrus.Infof("Handling %s", event.GetAction())
	switch event.GetAction() {
	case "opened":
		if err := h.opened(ctx, event); err != nil {
			return errors.Wrap(err, "failed to parse pr")
		}

	}

	return nil
}

func (h *PullRequestHandler) opened(ctx context.Context, event github.PullRequestEvent) error {
	installationID := githubapp.GetInstallationIDFromEvent(&event)
	client, err := h.NewInstallationClient(installationID)
	if err != nil {
		return err
	}

	repo := event.GetRepo()
	repoOwner := repo.GetOwner().GetName()
	repoName := repo.GetName()
	author := event.GetSender()

	issueNumber, err := FindIssueNumberByAssignee(ctx, client, repoOwner, repoName, author.GetLogin())
	if err != nil {
		return err
	} else if issueNumber == 0 {
		return nil
	}

	comment := github.PullRequestComment{
		Body: String(fmt.Sprintf(`## Step 5: Link a Pull Request to an Issue

Awesome work creating that PR.  

Now let's link it to our issue so that when the PR is merged, GitHub will automatically resolve our Issue.

### :keyboard: Activity: Edit a pull request

1. Click on the **...** icon located at the top right corner of the first comment's box, then click on **Edit** to make an edit
1. Add a description of the changes you've made in the comment box. Feel free to add a description of what you’ve accomplished so far. As a reminder, you have: created a branch, created a file and made a commit, and opened a pull request
1. Add the text "Resolves #%d" to link this PR with that Issue.
1. Click the green **Update comment** button at the bottom right of the comment box when done

<hr>
<h3 align="center">I'll respond when I detect this pull request's body has been edited.</h3>`, issueNumber)),
	}
	if _, _, err := client.PullRequests.CreateComment(ctx, repoOwner, repoName, event.GetPullRequest().GetNumber(), &comment); err != nil {
		logrus.WithError(err).Error("Failed to create pr comment")
	}

	return nil
}
