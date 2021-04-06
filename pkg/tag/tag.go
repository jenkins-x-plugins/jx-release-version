package tag

import (
	"fmt"
	"os"

	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	"github.com/go-git/go-git/v5/config"

	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
)

type Tag struct {
	FormattedVersion string
	Dir              string
	PushTag          bool
}

func (options Tag) TagRemote() error {
	var err error
	if options.Dir == "" {
		options.Dir, err = os.Getwd()
		if err != nil {
			return errors.Wrapf(err, "failed to get current directory")
		}
	}
	if options.FormattedVersion == "" {
		return errors.Wrapf(err, "no version to use for tag")
	}

	repo, err := git.PlainOpen(options.Dir)
	if err != nil {
		return errors.Wrapf(err, "failed to open dir %s", options.Dir)
	}

	h, err := repo.Head()
	if err != nil {
		return errors.Wrap(err, "failed to get HEAD commit")
	}

	tagOptions := &git.CreateTagOptions{
		Message: fmt.Sprintf("Release version %s", options.FormattedVersion),
	}
	log.Logger().Debugf("git tag -a %s -m \"%s\"", options.FormattedVersion, tagOptions.Message)
	_, err = repo.CreateTag(options.FormattedVersion, h.Hash(), tagOptions)
	if err != nil {
		return errors.Wrap(err, "failed to create tag")
	}

	if options.PushTag {
		return pushTags(repo)
	}
	return nil
}

func pushTags(r *git.Repository) error {

	po := &git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
	}
	log.Logger().Debug("git push --tags")
	err := r.Push(po)

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Logger().Debug("origin remote was up to date, no push done")
			return nil
		}
		return err
	}
	return nil
}
