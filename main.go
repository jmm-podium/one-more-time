package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

var (
	githubTokenFlag string
	githubRepoFlag  string
	gitCommitFlag   string
)

func main() {
	ctx := context.Background()
	_, _ = fmt.Fprintf(os.Stdout, "the token is set? %v\n", githubTokenFlag != "")
	_, _ = fmt.Fprintf(os.Stdout, "github repo? %s\n", githubRepoFlag)
	_, _ = fmt.Fprintf(os.Stdout, "git commit %s\n", gitCommitFlag)

	githubOrg, githubProject := parseRepo(githubRepoFlag)

	githubClient := makeGithubClient(ctx, githubTokenFlag)
	pullRequests, _, err := githubClient.PullRequests.List(
		ctx,
		githubOrg,
		githubProject,
		&github.PullRequestListOptions{
			State: "closed",
		})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var mergedPullRequests []github.PullRequest
	for _, pullRequest := range pullRequests {
		if pullRequest == nil || pullRequest.MergedAt == nil {
			continue
		}
		mergedPullRequests = append(mergedPullRequests, *pullRequest)
	}

	sortPullRequests(mergedPullRequests)

	bytes, _ := json.MarshalIndent(mergedPullRequests, "", "  ")
	_, _ = fmt.Fprint(os.Stdout, string(bytes)+"\n")

	result := "great success"
	_, _ = fmt.Fprintf(os.Stdout, "::set-output name=result::%s\n", result)
}

func init() {
	flag.StringVar(&githubTokenFlag, "github-token", "", "GitHub Access Token for accessing the wiki repo")
	flag.StringVar(&githubRepoFlag, "github-repo", "", "GitHub repository")
	flag.StringVar(&gitCommitFlag, "git-commit", "", "git commit hash")
	flag.Parse()
}

func makeGithubClient(ctx context.Context, token string) *github.Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tokenClient := oauth2.NewClient(ctx, tokenSource)
	return github.NewClient(tokenClient)
}
func parseRepo(githubRepo string) (string, string) {
	githubRepoParts := strings.Split(githubRepo, "/")
	return githubRepoParts[0], githubRepoParts[1]
}
func sortPullRequests(pullRequests []github.PullRequest) {
	sort.Slice(pullRequests, func(this, that int) bool {
		thisMergedAt := *pullRequests[this].MergedAt
		thatMergedAt := *pullRequests[that].MergedAt
		return thisMergedAt.After(thatMergedAt)
	})
}
