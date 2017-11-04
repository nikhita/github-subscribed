package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

const (
	// BANNER is what is printed for help/info output.
	BANNER = "github-subscribed : %s\n\n"
	// USAGE is an example of how the command should be used.
	USAGE = "USAGE:\ngithub-subscribed -token=<your-token> -repo=<repo>"
	// VERSION is the binary version.
	VERSION = "v0.1.0"
)

var (
	token   string
	version bool
	repo    string
)

func init() {
	flag.StringVar(&token, "token", "", "Mandatory GitHub API token")
	flag.StringVar(&repo, "repo", "", "(optional) Search threads belonging to a particular repo.\n\tYou must own the repo or be a member of the organization which owns the repo.")
	flag.BoolVar(&version, "version", false, "Print version and exit")
	flag.BoolVar(&version, "v", false, "(shorthand) Print version and exit")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, VERSION))
		fmt.Println(USAGE)
		flag.PrintDefaults()
	}

	flag.Parse()

	if version {
		fmt.Printf("%s", VERSION)
		os.Exit(0)
	}

	if token == "" {
		usageAndExit("GitHub token cannot be empty", 1)
	}
}

func main() {
	ctx := context.Background()

	// Create an authenticated client.
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	if len(repo) == 0 {
		fmt.Printf("Fetching threads you are subscribed to...\n")
	} else {
		fmt.Printf("Fetching threads you are subscribed in %s...\n", repo)
	}

	output := getSubscribedThreads(ctx, client, repo)

	fmt.Println("Total number of subscribed threads: ", len(output))
	for num, line := range output {
		serialNumber := fmt.Sprintf("%v.", num+1)
		fmt.Printf("%s %s\n", serialNumber, line)
	}
}

func getSubscribedThreads(ctx context.Context, client *github.Client, repo string) []string {
	var subscribedIssues []string

	opt := github.IssueListOptions{
		Filter:      "subscribed",
		ListOptions: github.ListOptions{PerPage: 50},
	}

	for {
		// rate limit is 5000 requests per hour, so probably won't happen.
		sleepIfRateLimitExceeded(ctx, client)

		issuesResults, resp, err := client.Issues.List(ctx, true, &opt)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, issue := range issuesResults {
			if len(repo) != 0 && issue.Repository.GetName() != repo {
				continue
			}

			formattedIssue := printInMarkdownFormat(issue)
			subscribedIssues = append(subscribedIssues, formattedIssue)
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return subscribedIssues
}

func printInMarkdownFormat(issue *github.Issue) string {
	issueLink := fmt.Sprintf("[#%v](%s) - ", issue.GetNumber(), issue.GetHTMLURL())
	issueTitle := fmt.Sprintf("%s", issue.GetTitle())
	// [#1234](<link>) - <title>
	formatedIssue := fmt.Sprintf("%s%s", issueLink, issueTitle)
	return formatedIssue
}

func sleepIfRateLimitExceeded(ctx context.Context, client *github.Client) {
	rateLimit, _, err := client.RateLimits(ctx)
	if err != nil {
		fmt.Printf("Problem in getting rate limit information %v\n", err)
		return
	}

	if rateLimit.Search.Remaining == 1 {
		timeToSleep := rateLimit.Search.Reset.Sub(time.Now()) + time.Second
		time.Sleep(timeToSleep)
	}
}

func usageAndExit(message string, exitCode int) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(exitCode)
}
