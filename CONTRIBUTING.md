# Contributing
We're happy that you want to contribute to this project. Please read the sections to make the process as smooth as possible.

## Requirements
- Golang `1.23`
- An OpenAI API key
  * OpenAI API keys can be obtained from [OpenAI](https://platform.openai.com/account/api-keys)
  * You can set the API key for Ava using `export OPENAI_API_KEY=$MY_API_KEY` or `.envrc`
- If you want to build the container image, you need to have a container engine (docker, podman, rancher, etc.) installed

## Getting Started

**Where should I start?**
- If you are looking for something to work on, check out our [open issues](https://github.com/matthisholleville/ava/issues).
- If you have an idea for a new feature, please open an issue, and we can discuss it.
- We are also happy to help you find something to work on. Just reach out to us.

**Discuss issues**
* Before you start working on something, propose and discuss your solution on the issue
* If you are unsure about something, ask the community

**How do I contribute?**
- Fork the repository and clone it locally
- Create a new branch and follow [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) guidelines for work undertaken
- Assign yourself to the issue, if you are working on it (if you are not a member of the organization, please leave a comment on the issue)
- Make your changes
- Keep pull requests small and focused, if you have multiple changes, please open multiple PRs
- Create a pull request back to the upstream repository and follow the [pull request template](.github/pull_request_template.md) guidelines.
- Wait for a review and address any comments

**Opening PRs**
- As long as you are working on your PR, please mark it as a draft
- Please make sure that your PR is up-to-date with the latest changes in `main`
- Fill out the PR template
- Mention the issue that your PR is addressing (closes: #<id>)
- Make sure that your PR passes all checks

**Reviewing PRs**
- Be respectful and constructive
- Assign yourself to the PR
- Check if all checks are passing
- Suggest changes instead of simply commenting on found issues
- If you are unsure about something, ask the author
- If you are not sure if the changes work, try them out
- Reach out to other reviewers if you are unsure about something
- If you are happy with the changes, approve the PR
- Merge the PR once it has all approvals and the checks are passing

## Semantic commits
We use [Semantic Commits](https://www.conventionalcommits.org/en/v1.0.0/) to make it easier to understand what a commit does and to build pretty changelogs. Please use the following prefixes for your commits:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `chores`: Changes to the build process or auxiliary tools and libraries such as documentation generation
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `test`: Adding missing tests or correcting existing tests
- `ci`: Changes to our CI configuration files and scripts

An example for this could be:
```
git commit -m "docs: add a new section to the README"
```

## Building
Building the binary is as simple as running `make build` in the root of the repository. If you want to build the container image, you can run `make build-container` in the root of the repository.

## Releasing
Releases of Ava are done using [Release Please](https://github.com/googleapis/release-please) and [GoReleaser](https://goreleaser.com/). The workflow looks like this:

* A PR is merged to the `main` branch:
  * Release please is triggered, creates or updates a new release PR
  * This is done with every merge to main, the current release PR is updated every time

* Merging the 'release please' PR to `main`:
  * Release please is triggered, creates a new release and updates the changelog based on the commit messages
  * GoReleaser is triggered, builds the binaries and attaches them to the release
  * Containers are created and pushed to the container registry

> With the next relevant merge, a new release PR will be created and the process starts again

