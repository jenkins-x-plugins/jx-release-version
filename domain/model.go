package domain

import "context"

type GitClient interface {
	ListTags(ctx context.Context, owner string, repo string) ([]Tag, error)
}

type Tag struct {
	Name string
}
