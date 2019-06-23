# Git Training

This is a GitHub App that you apply to a training repo, and it interacts with trainees through the process of using Issues, PRs, creating files, branches, and resolving.

# Process

1. User opens Issue
2. Bot comments on issue asking for them to self-assign

- Hook: issue create
- Validate: none
- Action: make comment on issue

3. User assigns issue to themself
4. Bot comments on issue asking user to create a new file on a branch named username-branch-1 containing "First commit!", then create a pull request

- Hook: issue assigned
- Validate: issue assignee matches author
- Action: make comment on issue

5. User creates branch and PR
6. Bot sees new PR & reviews it asking for a change to associate it with the issue

- Hook: pr opened
- Validate: contents of file, name of branch
- Action: make comment on pr

7. User makes change which updates branch/pr
8. Bot approves PR and asks user to merge

- Hook: pr updated
- Validate: contents of file,
- Action: approve pr, make comment on pr

9. User merges PR which closes original issue
10. Bot asks user to verify the original issue was resolved

- Hook: pr closed
- Validate: pr & issue both merged
- Action: make comment on pr
