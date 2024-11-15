package tag

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
)

type Tag struct {
	FormattedVersion string
	Dir              string
	PushTag          bool
	GitName          string
	GitEmail         string
}

func (options Tag) TagRemote() error {
	var err error
	if options.Dir == "" {
		options.Dir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}
	}
	if options.FormattedVersion == "" {
		return errors.New("no version to use for tag")
	}

	repo, err := git.PlainOpen(options.Dir)
	if err != nil {
		return fmt.Errorf("failed to open git repository at %q: %w", options.Dir, err)
	}

	h, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get the HEAD reference for git repository at %q: %w", options.Dir, err)
	}

	tagOptions := &git.CreateTagOptions{
		Message: fmt.Sprintf("Release version %s", options.FormattedVersion),
	}

	// override default git config tagger info
	if options.GitName != "" && options.GitEmail != "" {
		log.Logger().Debugf("overriding default git tagger config with name %s: email: %s", options.GitName, options.GitEmail)
		tagOptions.Tagger = &object.Signature{
			Name:  options.GitName,
			Email: options.GitEmail,
			When:  time.Now(),
		}
	}

	log.Logger().Debugf("git tag -a %s -m %q", options.FormattedVersion, tagOptions.Message)
	_, err = repo.CreateTag(options.FormattedVersion, h.Hash(), tagOptions)
	if err != nil {
		return fmt.Errorf("failed to create tag %q with message %q: %w", options.FormattedVersion, tagOptions.Message, err)
	}

	if options.PushTag {
		return pushTags(repo)
	}
	return nil
}

func pushTags(r *git.Repository) error {
	token := os.Getenv("GIT_TOKEN")

	po := &git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stderr,
		RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
	}

	if token != "" {
		user := os.Getenv("GIT_USER")
		if user == "" {
			user = "abc123" // yes, this can be anything except an empty string
		}
		po.Auth = &http.BasicAuth{
			Username: user,
			Password: token,
		}
	}
	log.Logger().Debug("git push --tags")
	err := r.Push(po)

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Logger().Debug("origin remote was up to date, no push done")
			return nil
		}
		return fmt.Errorf("failed to push tags to origin: %w", err)
	}
	return nil
}
