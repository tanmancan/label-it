# Label It

A command line tool for adding labels to Github pull requests using a declarative YAML configuration. Written in Go.

## Download

Available for Mac, Linux or Windows. Latest release can be found [here](https://github.com/tanmancan/label-it/releases/latest).

Download and extract the `tar.gz` file.

```bash
 tar -xzf [FILENAME].tar.gz
```

Verify the tool is working:

```bash
./label-it --version

Version: vX.X.X
API Version: vX
SHA: XXXXXX
```

Scroll down to learn more about configuring rules and usage.

### Manual build

There is a [Makefile](Makefile) included with this project, which you can use to compile and build the project yourself. This includes tests and targets for Mac, linux and windows 64bit platforms. Binaries will be output in the `bin` directory.

## Usage

Configure Rules via an YAML file. See below for all [configuration options](#configuration-options)

`label-it.yaml`:
```yaml
apiVersion: 1
access:
  user: $GIT_USER
  token: $GIT_TOKEN
owner: tanmancan
repo: label-it
rules:
  - label: My Label Name
    head-rule:
      exact: head-branch-name
...
```

Run `label-it` and pass the configuration file as an option.
```bash
./label-it -c label-it.yaml
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
Version of rule configuration schema. Breaking changes to the schema will be versioned.

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

It is recommended that you pass in the authentication information via an env variable. Values that begin with a `$` will be treated as an env variable:

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

## Rule Types

### `base-rule`

Rule type that compares the pull request base branch.

```yaml
base-rule:
  match: master
  no-match: ^(dev-)

```

### `head-rule`

Rule type that compares the title text of the pull request.

```yaml
head-rule:
  exact: master
```

### `title-rule`
Rule type that compares the title text of the pull request.

```yaml
title-rule:
  no-match: My pull request title
```


### `body-rule`
Rule type that compares the body text of the pull request.

```yaml
body-rule:
  match: My pull request body
```

### `user-rule`
Rule type that compares the username of the account that opened the pull request.

```yaml
user-rule:
  no-exact: tanmancan
```

### `number-rule`
Rule type that compares the pull request number.

```yaml
number-rule:
  match: ^(10)
  no-exact: 1044
```

### `file-rule`
Rule type that compares file path of all changed files in a pull request. Currently only supports a maximum of 1000 files. If your pull request is larger, it is recommended manually adding labels for now. In the future a larger file count may be supported.

```yaml
file-rule:
  exact: www/index.html
  no-exact: files/common/robots.text
  match: (.jpg)$
  no-match: ^(readme.)
```

Since this rule compares a list of file paths, each check will need to validate against the full list:
- `exact`: If an exact match is found in the list of file paths,  exact check will be considered valid
- `no-exact`: If a no exact match is found in the list of file paths, no-exact check will be considered invalid.
- `match`: If an regex pattern matches a path in the list of file paths, then the match check will be considered valid.
- `no-match`: If an regex pattern matches a path in the list of file paths, then the no-match check will be considered invalid.

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

    # Rule type that compares the pull request head branch.
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

    # Rule type that compares the pull request base branch.
    # Each rule type may have four checks: exact, no-exact, match, no-match.
    base-rule:
      exact: master
      no-exact: staging
      match: ^(mas)
      no-match: ^(stag)

    # Rule type that compares the title text of the pull request.
    # Each rule type may have four checks: exact, no-exact, match, no-match.
    title-rule:
      exact: master
      no-exact: staging
      match: ^(mas)
      no-match: ^(stag)

    # Rule type that compares the body text of the pull request.
    # Each rule type may have four checks: exact, no-exact, match, no-match.
    body-rule:
      exact: master
      no-exact: staging
      match: ^(mas)
      no-match: ^(stag)

    # Rule type that compares the username of the account that opened the pull request.
    # Each rule type may have four checks: exact, no-exact, match, no-match.
    user-rule:
      exact: tanmancan
      no-exact: octocat
      match: ^(tan)
      no-match: ^(linus)

    # Rule type that compares the pull request number.
    # Each rule type may have four checks: exact, no-exact, match, no-match.
    number-rule:
      exact: 42
      no-exact: 10
      match: (100)|(200)|^(3)
      no-match: ^(5)|(600)

    # Rule type that compares paths of all changed files in a pr
    # Currently only supports testing the first 1000 files in a pr.
    # If your pull request is larger, it is recommended manually adding labels for now. In the future a larger file count may be supported.
    # Each rule type may have four checks: exact, no-exact, match, no-match.
    file-rule:
      # If an exact match is found in the list of file paths,
      # exact check will be considered valid
      exact: www/index.html
      # If a no exact match is found in the list of file paths,
      # no-exact check will be considered invalid
      no-exact: files/common/robots.text
      # If an regex pattern matches a path in the list of file paths,
      # then the match check will be considered valid
      match: (.jpg)$
      # If an regex pattern matches a path in the list of file paths,
      # then the no-match check will be considered invalid
      no-match: ^(readme.)

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

    # The label Robots will be applied if the file
    # www/robots.txt has been modified din the pull request
  - label: Robots
    file-rule:
      exact: www/robots.txt
```

---

&copy; Tanveer Karim [www.tkarimdesign.com](https://www.tkarimdesign.com)
