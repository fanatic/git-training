package handlers

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

func FindIssueNumberByAssignee(ctx context.Context, client *github.Client, repoOwner, repoName, assignee string) (int, error) {
	issues, _, err := client.Issues.ListByRepo(ctx, repoOwner, repoName, &github.IssueListByRepoOptions{
		Assignee: assignee,
	})
	if err != nil {
		return 0, err
	}

	if len(issues) == 0 {
		logrus.Infof("Dropping created event because no issues in repo assigned to %s", assignee)
		return 0, nil
	}
	return issues[0].GetNumber(), nil
}
