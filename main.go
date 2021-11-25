package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func listOrgsRepos(orgname string, apikey string) []*github.Repository {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apikey},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repos, _, _ := client.Repositories.ListByOrg(context.Background(), orgname, nil)

	return repos
}

func listUsersRepos(username string, apikey string) []*github.Repository {
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

func listCommits(ownername string, reponame string, apikey string) []string {
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

func usage() {
	fmt.Println("usage: github2email <UserName or OrgName>")
}

func main() {
	ownername := os.Args[1]
	apikey := os.Getenv("GITHUB_TOKEN")
	if ownername == "" {
		usage()
		os.Exit(0)
	} else if apikey == "" {
		fmt.Println("Set GitHub Access Token as GITHUB_TOKEN")
		os.Exit(0)
	}

	authors := []string{}
	repos := listUsersRepos(ownername, apikey)
	if len(repos) == 0 {
		repos = listOrgsRepos(ownername, apikey)
	}
	if len(repos) == 0 {
		fmt.Println("No public repository was found.")
		os.Exit(0)
	}

	for _, repo := range repos {
		authors = append(authors, listCommits(ownername, *repo.Name, apikey)...)
	}

	for _, author := range unique(authors) {
		fmt.Println(author)
	}
}
