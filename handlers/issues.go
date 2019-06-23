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
		logrus.Infof("Handling %s", event.GetAction())
		if err := h.opened(ctx, event); err != nil {
			return errors.Wrap(err, "failed to parse issue open")
		}
	case "assigned":
		logrus.Infof("Handling %s", event.GetAction())
		if err := h.assigned(ctx, event); err != nil {
			return errors.Wrap(err, "failed to parse issue open")
		}
	default:
		logrus.Infof("Handling %s", event.GetAction())
	}

	return nil
}

func (h *IssuesHandler) opened(ctx context.Context, event github.IssuesEvent) error {
	installationID := githubapp.GetInstallationIDFromEvent(&event)
	client, err := h.NewInstallationClient(installationID)
	if err != nil {
		return err
	}

	repo := event.GetRepo()
	issueNumber := event.GetIssue().GetNumber()
	repoOwner := repo.GetOwner().GetLogin()
	repoName := repo.GetName()
	author := event.GetIssue().GetUser()
	comment := github.IssueComment{
		Body: String(fmt.Sprintf(`# :wave: Welcome to GitHub Training, @%s!

I’ll guide you through some important first steps in coding and collaborating on GitHub.

This is an issue <sup>[:book:](https://help.github.com/articles/github-glossary/#issue)</sup>: a place to record bugs, request enhancements, or answer questions about your repo.

<hr>
<h3 align="center">Keep reading below to find your first task</h3>`, author.GetLogin())),
	}
	if _, _, err := client.Issues.CreateComment(ctx, repoOwner, repoName, issueNumber, &comment); err != nil {
		logrus.WithError(err).Error("Failed to create issue comment")
	}

	comment = github.IssueComment{
		Body: String(fmt.Sprintf(`## Step 1: Assign yourself

Unassigned issues don't have owners to look after them. 

### :keyboard: Activity

1. On the right side of the screen, under the "Assignees" section, click the gear icon and select yourself
		
<hr>
<h3 align="center">I'll respond when I detect you've assigned yourself to this issue.</h3>

> If you perform an expected action and don't see a response from me, wait a few seconds and refresh the page for your next steps.`)),
	}
	if _, _, err := client.Issues.CreateComment(ctx, repoOwner, repoName, issueNumber, &comment); err != nil {
		logrus.WithError(err).Error("Failed to create issue comment 2")
	}

	return nil
}

func String(s string) *string {
	return &s
}

func (h *IssuesHandler) assigned(ctx context.Context, event github.IssuesEvent) error {
	installationID := githubapp.GetInstallationIDFromEvent(&event)
	client, err := h.NewInstallationClient(installationID)
	if err != nil {
		return err
	}

	repo := event.GetRepo()
	issueNumber := event.GetIssue().GetNumber()
	repoOwner := repo.GetOwner().GetLogin()
	repoName := repo.GetName()
	author := event.GetIssue().GetUser()
	assignee := event.GetIssue().GetAssignee()

	if author.GetLogin() != assignee.GetLogin() {
		return nil
	}

	comment := github.IssueComment{
		Body: String(fmt.Sprintf(`## Introduction to a typical workflow

Now that you're familiar with issues, let's use this issue to track your path to your first contribution.

People use different workflows to contribute to software projects, but the simplest and most effective way to contribute on GitHub is the GitHub flow.

:tv: [Video: Understanding the GitHub flow](https://www.youtube.com/watch?v=PBI2Rz-ZOxU)

<hr>
<h3 align="center">Read below for next steps</h3>`)),
	}
	if _, _, err := client.Issues.CreateComment(ctx, repoOwner, repoName, issueNumber, &comment); err != nil {
		logrus.WithError(err).Error("Failed to create issue comment")
	}

	comment = github.IssueComment{
		Body: String(fmt.Sprintf(`## Step 2: Create a branch

Let’s complete the first step of the GitHub flow: creating a branch <sup>[:book:](https://help.github.com/articles/github-glossary/#branch)</sup>.

### :keyboard: Activity: Your first branch

1. Navigate to the [Code tab](https://github.factset.com/%s/%s)
2. Click **Branch: master** in the drop-down
3. In the field, enter a name for your branch, like "feat/username-1"
4. Click **Create branch: <name>** or press the “Enter” key to create your branch


<hr>
<h3 align="center">I'll respond when I detect a new branch has been created in this repository.</h3>`, repoOwner, repoName)),
	}
	if _, _, err := client.Issues.CreateComment(ctx, repoOwner, repoName, issueNumber, &comment); err != nil {
		logrus.WithError(err).Error("Failed to create issue comment 2")
	}

	return nil
}
