# Git Training

This is a GitHub App that you apply to a training repo, and it interacts with trainees through the process of using Issues, PRs, creating files, branches, and resolving.

# Process

Master branch is protected & no PR without 1 approving review

1. User opens Issue
2. Bot comments on issue asking for them to self-assign

- Hook: issue create
- Validate: none
- Action: make comment on issue

3. User assigns issue to themself
4. Bot comments on issue asking user to create a new branch named feat/username-1

- Hook: issue assigned
- Validate: issue assignee matches author
- Action: make comment on issue

5. User creates branch
6. Bot comments on issue asking user to create a new file on the branch called "First commit!", then create a pull request

- Hook: branch created
- Validate: branch name matches creator
- Action: make comment on issue

7. User creates file
8. Bot comments on issue asking user to create a pull request

- Hook: push
- Validate: file created
- Action: make comment on issue

9. User creates PR
10. Bot sees new PR & reviews it asking for a change to associate it with the issue (add resolves #1 to body of PR)

- Hook: pr opened
- Validate: contents of file, name of branch
- Action: make comment on pr

11. User makes change which updates pr
12. Bot adds review requesting update to the file

- Hook: pr updated
- Validate: contents of file,
- Action: add review

13. User makes change which updates branch
14. Bot approves PR and asks user to merge

- Hook: branch updated
- Validate: contents of file,
- Action: approve pr, make comment on pr

15. User merges PR which closes original issue
16. Bot asks user to verify the original issue was resolved

- Hook: pr closed
- Validate: pr & issue both merged
- Action: make comment on pr
