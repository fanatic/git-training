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

type CreateHandler struct {
	githubapp.ClientCreator
}

func (h *CreateHandler) Handles() []string {
	return []string{"create"}
}

func (h *CreateHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.CreateEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse create event payload")
	}

	switch event.GetRefType() {
	case "branch":
		logrus.Infof("Handling %s", event.GetRefType())
		if err := h.branchCreated(ctx, event); err != nil {
			return errors.Wrap(err, "failed to parse create")
		}
	default:
		logrus.Infof("Handling %s", event.GetRefType())
	}

	return nil
}

func (h *CreateHandler) branchCreated(ctx context.Context, event github.CreateEvent) error {
	installationID := githubapp.GetInstallationIDFromEvent(&event)
	client, err := h.NewInstallationClient(installationID)
	if err != nil {
		return err
	}

	repo := event.GetRepo()
	repoOwner := repo.GetOwner().GetLogin()
	repoName := repo.GetName()
	author := event.GetSender()
	branchName := event.GetRef()

	issues, _, err := client.Issues.ListByRepo(ctx, repoOwner, repoName, &github.IssueListByRepoOptions{
		Assignee: author.GetLogin(),
	})
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		logrus.Infof("Dropping created event because no issues in repo assigned to %s", author.GetLogin())
		return nil
	}
	issueNumber := issues[0].GetNumber()

	comment := github.IssueComment{
		Body: String(fmt.Sprintf(`## Step 5: Commit a file

:tada: You created a branch!

Creating a branch allows you to make modifications to your project without changing the deployed "master" branch. Now that you have a branch, it’s time to create a file and make your first commit!

Commits are snapshots of file changes, so let's make our first one.

### :keyboard: Activity: Your first commit

1. Create a new file on this branch named with your username.
			- Return to the "Code" tab
			- In the branch drop-down, select "%s"
			- Click **Create new file**
			- In the "file name" field, type "users/%s.md". Entering the "/" in the filename will automatically place your file in the "users" directory.
1. When you’re done naming the file, add the following content to your file:
			`+"```yaml\n"+
			"Hello, world!\n"+
			"```"+`
1. After adding the text, you can commit the change by entering a commit message in the text-entry field below the file edit view.
1. When you’ve entered a commit message, click **Commit new file**

<hr>
<h3 align="center">I'll respond when I detect a new commit on this branch.</h3>`, branchName, author.GetLogin())),
	}
	if _, _, err := client.Issues.CreateComment(ctx, repoOwner, repoName, issueNumber, &comment); err != nil {
		logrus.WithError(err).Error("Failed to create issue comment")
	}

	return nil
}
