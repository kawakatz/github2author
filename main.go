package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func list_orgs_repo(orgname string, apikey string) []*github.Repository {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apikey},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repos, _, _ := client.Repositories.ListByOrg(context.Background(), orgname, nil)

	return repos
}

func list_users_repo(username string, apikey string) []*github.Repository {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apikey},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repos, _, _ := client.Repositories.List(context.Background(), username, nil)

	return repos
}

func unique(arr []string) []string {
	occurred := map[string]bool{}
	result := []string{}
	for e := range arr {

		// check if already the mapped
		// variable is set to true or not
		if !occurred[arr[e]] {
			occurred[arr[e]] = true

			// Append to result slice.
			result = append(result, arr[e])
		}
	}

	return result
}

func list_commits(ownername string, reponame string, apikey string) []string {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apikey},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	authors := []string{}
	//opt := &github.CommitsListOptions{Author: ownername}
	commitsInfo, _, _ := client.Repositories.ListCommits(context.Background(), ownername, reponame, nil)
	for _, commitInfo := range commitsInfo {
		authors = append(authors, *commitInfo.Commit.Author.Name+" <"+*commitInfo.Commit.Author.Email+">")
	}

	return authors
}

func main() {
	ownername := os.Args[1]
	apikey := os.Getenv("GITHUB_ACCESS_TOKEN")
	if ownername == "" {
		fmt.Println("Usage: github2email <UserName or OrgName>")
		os.Exit(0)
	} else if apikey == "" {
		fmt.Println("Set GitHub Access Token as GITHUB_ACCESS_TOKEN")
		os.Exit(0)
	}

	authors := []string{}
	repos := list_users_repo(ownername, apikey)
	if len(repos) == 0 {
		repos = list_orgs_repo(ownername, apikey)
	}
	if len(repos) == 0 {
		fmt.Println("No public repository was found.")
		os.Exit(0)
	}

	for _, repo := range repos {
		authors = append(authors, list_commits(ownername, *repo.Name, apikey)...)
	}

	for _, author := range unique(authors) {
		fmt.Println(author)
	}
}
