package main

import (
	"fmt"
	"os"
)

type ReleaseInfo struct {
	ID      uint
	TagName string
}

type ReleaseInfoer interface {
	GetLatestReleaseTag(string) (string, error)
}

type GitHubReleaseInfoer struct {
}

func (g GitHubReleaseInfoer) GetLatestReleaseTag(repo string) (tag string, err error) {
	return "", nil
}

func main() {
	g := GitHubReleaseInfoer{}
	msg, err := getReleaseTagMessage(g, "docker/machine")
	if err != nil {
		fmt.Fprintln(os.Stderr, msg)
	}
	fmt.Println(msg)
}

func getReleaseTagMessage(releaseInfoer ReleaseInfoer, repo string) (msg string, err error) {
	tag, err := releaseInfoer.GetLatestReleaseTag(repo)
	if err != nil {
		return "", fmt.Errorf("error querying Github API: %s", err)
	}
	return fmt.Sprintf("The latest release is %s", tag), nil
}
