package adapters

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/jenkins-x/jx-release-version/domain"
	"golang.org/x/oauth2"
	"os"
)

type GitHubClient struct {
	client *github.Client
}

func (g *GitHubClient) ListTags(ctx context.Context, owner string, repo string) ([]domain.Tag, error) {
	tags, _, err := g.client.Repositories.ListTags(ctx, owner, repo, nil)

	var a []domain.Tag

	for _, e := range tags {
		a = append(a, domain.Tag{Name: e.GetName()})
	}

	return a, err
}

func NewGitHubClient(ghOwner string, ghRepository string, debug bool) domain.GitClient {
	var githubClient domain.GitClient

	if ghOwner != "" && ghRepository != "" {
		token := os.Getenv("GITHUB_AUTH_TOKEN")
		ctx := context.Background()
		if token != "" {
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			)
			tc := oauth2.NewClient(ctx, ts)

			githubClient = &GitHubClient{
				client: github.NewClient(tc),
			}
		} else {
			if debug {
				fmt.Println("no GITHUB_AUTH_TOKEN env var found so using unauthenticated request")
			}
			githubClient = &GitHubClient{
				client: github.NewClient(nil),
			}
		}
	}

	return githubClient
}
