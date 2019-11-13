package adapters

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/jenkins-x/jx-release-version/domain"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

type GitHubClient struct {
	client *github.Client
}

func (g *GitHubClient) ListTags(ctx context.Context, owner string, repo string) ([]domain.Tag, error) {
	tags, _, err := g.client.Repositories.ListTags(ctx, owner, repo, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting tags: %v", err)
	}

	var a []domain.Tag

	for _, e := range tags {
		a = append(a, domain.Tag{Name: e.GetName()})
	}

	return a, err
}

func NewGitHubClient(debug bool) domain.GitClient {
	var oauth2Client *http.Client

	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		oauth2Client = oauth2.NewClient(context.Background(), ts)
	} else {
		if debug {
			fmt.Println("no GITHUB_AUTH_TOKEN env var found so using unauthenticated request")
		}
	}

	return &GitHubClient{
		client: github.NewClient(oauth2Client),
	}
}
