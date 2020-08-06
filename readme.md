# Label It

A command line tool written in Go, for adding labels to Github pull requests using a declarative YAML configuration

## Usage

Configure Rules via YAML. See below for all [configuration options](#options)

`label-it.yaml`:
```yaml
apiVersion: 1
access:
  user: $GIT_USER
  token: $GIT_TOKEN
owner: tanmancan
repo: label-it
rules:
  MyLabel:
    base-rule: dev
...
```

Run using rules from the YAML configuration
```bash
label-it -c label-it.yaml
```

## Usage Options

View available options via `label-it --help`
```
Usage: ./label-it [--version][--help][-c <path>]
Example: ./label-it -c label-it.yaml

  --help
        Display the help text
  -c string
        Path to the yaml file
  -dry
        Outputs list of pull request and matched labels. Does not call the API
  -version
        Show version information
```

## Configuration Options

### `apiVersion` (`int`) *required*
Version of YAML schema. Breaking changes to the schema will be versioned.

```yaml
apiVersion: 1
```

### `access` (`map`) *required*
Github username and personal access token. This will be used to authenticate with the API. Token will require the `repo` scope in order to view and update existing pull requests. For more information see: https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token

```yaml
access:
  user: tanmancan
  token: asdf1234
```

While you can add a username and token directly to the YAML file, it is recommended that you pass in an env variable instead. Env variables are determined by prefixing the name with `$`:

```yaml
access:
  user: $GIT_USER
  token: $GIT_TOKEN
```

### `owner` (`string`) *required*
The repository owner

```yaml
owner: tanmancan
```

### `repo` (`string`) *required*
The repository name

```yaml
repo: label-it
```

### `rules` (`map`) *required*
Provide a list of rules, that are grouped by labels. If all rules in a group match a pull request, then the label will be added to the PR.

```yaml
rules:
  - label: My Label Name
    head-rule:
      exact: head-branch-name
    base-rule:
      exact: base-branch-name
```

### `label` (`string`)

The label name to be applied

```yaml
rules:
  - label: Label Name
```

If a label is provided with no rules (ex: `Add This Label To All PR`), then it will be added to all open pull requests

```yaml
rules:
  - label: Add This Label To All PR
  - label: another-label
      base-rule:
        exact: base-branch-name
```

## Rule Checks
Each rule types may have 4 possible checks. You can provide one or more checks for a given rule type.

```yaml
base-rule:
  exact: text
  no-exact: text
  match: text
  no-match: text
```

### `exact` (`string`)
The rule value must be an exact match of the compare value. In the example below, a pull request merging to the `master` branch will pass this check.

```yaml
base-rule:
  exact: master
```

### `no-exact` (`string`)
The rule value must NOT be an exact match of the compare value. In the example below, any pull request merging to the `master` branch will NOT pass this check. This check uses the go `regexp` library which uses the RE2 syntax. Learn more [here](https://github.com/google/re2/wiki/Syntax).

```yaml
base-rule:
  no-exact: master
```

### `match` (`string`)
A regex pattern that must match a compare value. In the example below, any pull request that is merging to a branch name staring with `stage-` will pass this check. This check uses the go `regexp` library which uses the RE2 syntax. Learn more [here](https://github.com/google/re2/wiki/Syntax).
```yaml
base-rule:
  match: ^(stage-)
```

### `no-match` (`string`)
 a regex pattern that must NOT match a compare value. In the example below, any pull request that is merging to a branch name staring with `stage-` will NOT pass this check.
```yaml
base-rule:
  no-match: ^(stage-)
```

## Rules Types

### `base-rule`

Rules type that compares the pull request base branch.

```yaml
base-rule:
  match: master
  no-match: ^(dev-)

```

### `head-rule`

Rules type that compares the title text of the pull request.

```yaml
head-rule:
  exact: master
```

### `title-rule`
Rules type that compares the title text of the pull request.

```yaml
title-rule:
  no-match: My pull request title
```


### `body-rule`
Rules type that compares the body text of the pull request.

```yaml
body-rule:
  match: My pull request body
```

### `user-rule`
Rules type that compares the username of the account that opened the pull request.

```yaml
user-rule:
  no-exact: tanmancan
```

### `number-rule`
Rules type that compares the pull request number.

```yaml
number-rule:
  match: ^(10)
  no-exact: 1044
```

### Example Configuration

```yaml
apiVersion: 1

# Github username and access token
# Recommend passing in an env variable by prefixing it with $
access:
  user: $GIT_USER
  token: $GIT_TOKEN

# Repository owner
owner: tanmancan

# Repository name
repo: label-it

# Provide a list of rules, that are grouped by labels
# If all rules in a group match a pull request,
# then the label will be added to the PR.
rules:

    # Label name. If all provided rules match,
    # then this label will be added to the pull request.
  - label: my-label-name

    # Rules type that compares the pull request head branch.
    # Each rule type may have four checks: exact, no-exact, match, no-match.
    head-rule:

      # Each rule types may have 4 possible checks. You can
      # provide one or more checks for a given rule type:
      # Exact - the rule value must be an exact match of the compare value.
      exact: master
      # NoExact - a rule value must NOT be an exact match of the compare value.
      no-exact: staging
      # Match - a regex pattern that must match a compare value.
      match: ^(mas)
      # NoMatch - a regex pattern that must NOT match a compare value
      no-match: ^(stag)

    # Rules type that compares the pull request base branch.
    # Each rule type may have four checks: exact, no-exact, match, no-match.
    base-rule:
      exact: master
      no-exact: staging
      match: ^(mas)
      no-match: ^(stag)

    # Rules type that compares the title text of the pull request.
    # Each rule type may have four checks: exact, no-exact, match, no-match.
    title-rule:
      exact: master
      no-exact: staging
      match: ^(mas)
      no-match: ^(stag)

    # Rules type that compares the body text of the pull request.
    # Each rule type may have four checks: exact, no-exact, match, no-match.
    body-rule:
      exact: master
      no-exact: staging
      match: ^(mas)
      no-match: ^(stag)

    # Rules type that compares the username of the account that opened the pull request.
    # Each rule type may have four checks: exact, no-exact, match, no-match.
    user-rule:
      exact: tanmancan
      no-exact: octocat
      match: ^(tan)
      no-match: ^(linus)

    # Rules type that compares the pull request number.
    # Each rule type may have four checks: exact, no-exact, match, no-match.
    number-rule:
      exact: 42
      no-exact: 10
      match: (100)|(200)|^(3)
      no-match: ^(5)|(600)

    # Examples:

    # If not rules are given, label
    # will be added to all pull requests
  - label: test

    # All rules must match for label to be added
  - label: All Rules Much Match
    # The pr head branch must exactly be "staging"
    head-rule:
      exact: staging
    # The pr base branch must not be "dev"
    base-rule:
      no-exact: dev
    # The pr title must be "Create"
    title-rule:
      exact: Create

    # Base rule, exact - If a pull request's base branch
    # is "master", the the label "To Master" will be added
  - label: To Master
    base-rule:
      exact: master

    # Head rule, no-exact - The pull request's head branch must not
    # be "staging", in order for the rule to be applied
  - label: From Staging
    head-rule:
      no-exact: staging

    # Head rule, match - This label will be applied if the pull request
    # head branch name starts with "pr-"
  - label: PR-* branch
    head-rule:
      match: ^(pr-)

    # Title rule, no-match - Any pull request that have the words
    # "text", "in", or "title" in its title, will not have this label applied.
  - label: Regex Title
    title-rule:
      no-match: (text)|(in)|(title)

    # The label "PR #5" will be added to pull request #5
  - label: "PR #5"
    number-rule:
      exact: 5

    # The label will be added if the user is not "tanmancan".
  - label: Tanmancan PR
    user-rule:
      no-exact: tanmancan

    # The label "BodyText" will be applied if the pull request contains
    # only the text "Hello World"
  - label: BodyText
    body-rule:
      exact: Hello World
```
