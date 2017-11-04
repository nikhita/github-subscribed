# github-subscribed

github-subscribed is a tool to list all threads (issues and pull requests) you are subscribed to.
Please note that the subscribed threads will only be generated for repositories you own or for organizations
you are a member of.

The output will be in the markdown format. You can copy-paste it into a markdown file or a gist for easy reading.

## Installation

**Prerequisites**: Go version 1.7 or greater.

1. Get the code

```
$ go get github.com/nikhita/github-subscribed
```

2. Build

```
$ cd $GOPATH/src/github.com/nikhita/github-subscribed
$ go install
```

## Usage

To authenticate, you will need a Github API token. You can find more details about generating an API token [here](https://github.com/blog/1509-personal-api-tokens).

```
github-subscribed : v0.1.0

USAGE:
github-subscribed -token=<your-token> -repo=<repo>
  -repo string
    	(optional) Search threads belonging to a particular repo.
	You must own the repo or be a member of the organization which owns the repo.
  -token string
    	Mandatory GitHub API token
  -v	(shorthand) Print version and exit
  -version
    	Print version and exit
```

## License

github-contrib is licensed under the [MIT License](/LICENSE).
