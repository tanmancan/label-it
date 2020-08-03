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
    base: dev
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
Rules describe conditions that a pull request must match. If an open pull request matches all provided rules then this label will be applied to it. Each rule is grouped by the label name as the key.

```yaml
rules:
  LabelName:
    base: base-branch-name
    head: head-branch-name
    title: Text in title
    body: Text in body
    user: tanmancan
    number: [0, 1, 2, 3, 4]
```

### `LabelName` (`string`)

Key is the Label name to be added to an open pull request

```yaml
MyLabel:
  base: base-branch-name
  head: head-branch-name
  ...
```

If a label is provided with no rules (ex: `AddToAllPr`), then it will be added to all open pull requests

```yaml
AddToAllPr:
MyLabel:
    base: base-branch-name
```

## Rules Options

### `base` (`string`)

The "to" branch. Any pull request merging to this branch will have this label added.

```yaml
base: master
```

May also pass a regex pattern
```yaml
base: ^(stage-)
```

### `head` (`string`)

The "from" branch. Any pull request created from this branch will have this label added.

```yaml
head: master
```

May also pass a regex pattern
```yaml
head: ^(hotfix-)
```

### `title` (`string`)

Rule checks if provided text appears in pull request title.

```yaml
title: My pull request title
```

May also pass a regex pattern
```yaml
title: (hotfix)|(bugfix)|(regressions)
```

### `body` (`string`)
Rule checks if provided text appears in pull request body. May pass a string or a regular expression

```yaml
body: My pull request body
```

May also pass a regex pattern
```yaml
body: (hotfix)|(bugfix)|(regressions)
```

### `user` (`string`)

Provide a Github username. If pull request was opened by this user, the label will be added

```yaml
user: tanmancan
```

### `number` (`array`)

Apply label if a pull request number matches any of the provided numbers.

```yaml
number: [5, 10, 15, 20]
```

### Example Configuration

```yaml
apiVersion: 1
access:
  user: $GIT_USER
  token: $GIT_TOKEN
owner: tanmancan
repo: label-it
rules:
  Label Name:
    base: base-branch-name
    head: head-branch-name
    title: Text in title
    body: Text in body
    user: tanmancan
    number: [0, 1, 2, 3, 4]

  # Examples:
  # All rules must match for label to be added
  AllRulesMatch:
    head: staging
    base: dev
    title: Create

  # A label with no rule will be added to all pull requests
  LabelToAllPr:

  # In this example - If a pull request's base branch
  # is "staging", the the label "From Staging" will be added
  From Staging:
    head: staging

  # If the pull request is merging to master branch
  # then the label "To Master" will be applied
  To Master:
    base: master

  # This label will be applied if the head branch
  # name starts with "pr-"
  PR-* branch:
    head: ^(pr-)

  # The label "PR #5" will be added to pull request #5
  # if it is open
  "PR #5":
    number: [5]

  # Regular expression string - if pull request has the words
  # text, in, or title the label "Regex Title" will be applied
  Regex Title:
    title: (text)|(in)|(title)

  # The label will be added if the user "tanmancan" opened the pull request
  Tanmancan PR:
    user: tanmancan

  # The label "BodyText" will be applied if the pull request contains
  # the text "Hello World"
  BodyText:
    body: Hello World

```
