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

type PushHandler struct {
	githubapp.ClientCreator
}

func (h *PushHandler) Handles() []string {
	return []string{"push"}
}

func (h *PushHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.PushEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse create event payload")
	}

	logrus.Infof("Handling %s", event.GetRef())

	if event.GetCreated() || event.GetDeleted() {
		logrus.Infof("Dropping push event because it was a create or delete")
		return nil
	}
	if err := h.commitCreated(ctx, event); err != nil {
		return errors.Wrap(err, "failed to parse push")
	}

	return nil
}

func (h *PushHandler) commitCreated(ctx context.Context, event github.PushEvent) error {
	installationID := githubapp.GetInstallationIDFromEvent(&event)
	client, err := h.NewInstallationClient(installationID)
	if err != nil {
		return err
	}

	repo := event.GetRepo()
	repoOwner := repo.GetOwner().GetName()
	repoName := repo.GetName()
	author := event.GetSender()
	branchName := event.GetRef()

	issueNumber, err := FindIssueNumberByAssignee(ctx, client, repoOwner, repoName, author.GetLogin())
	if err != nil {
		return err
	} else if issueNumber == 0 {
		return nil
	}

	expectedFilename := "users/" + author.GetLogin() + ".md"

	hasExpectedFile := false
	for _, filename := range event.GetHeadCommit().Added {
		if filename == expectedFilename {
			hasExpectedFile = true
		}
	}
	if !hasExpectedFile {
		comment := github.IssueComment{
			Body: String(fmt.Sprintf(`## Something's not quite right.

I'm looking for a new file named "users/%s.md" in your branch %s.`, author.GetLogin(), branchName)),
		}
		if _, _, err := client.Issues.CreateComment(ctx, repoOwner, repoName, issueNumber, &comment); err != nil {
			logrus.WithError(err).Error("Failed to create issue comment")
		}

		return nil
	}

	comment := github.IssueComment{
		Body: String(fmt.Sprintf(`## Step 4: Open a pull request

Nice work making that commit :sparkles:

In the real world, that commit would contain code working towards some feature or bug fix for one of our products.  Since we're just training here, it can contain anything.

Now that you’ve created a commit, it’s time to share your proposed change through a pull request! Where issues encourage discussion with other contributors and collaborators on a project, pull requests help you share your changes, receive feedback on them, and iterate on them until they’re perfect!

### :keyboard: Activity: Create a pull request

1. Open a pull request:
		- From the "Pull requests" tab, click **New pull request**
		- In the "base:" drop-down menu, make sure the "master" branch is selected
		- In the "compare:" drop-down menu, select "%s"
1. When you’ve selected your branch, enter a title for your pull request. For example "Add %s's file"
1. The next field helps you provide a description of the changes you made. Feel free to add a description of what you’ve accomplished so far. As a reminder, you have: created a branch, created a file and made a commit, and opened a pull request
1. Click **Create pull request**

<hr>
<h3 align="center">I'll respond in your new pull request.</h3>
		`, branchName, author.GetLogin())),
	}
	if _, _, err := client.Issues.CreateComment(ctx, repoOwner, repoName, issueNumber, &comment); err != nil {
		logrus.WithError(err).Error("Failed to create issue comment")
	}

	return nil
}
