package tag

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
)

func TestTag(t *testing.T) {

	// setup test
	dir := t.TempDir()

	r, err := git.PlainInit(dir, false)
	assert.NoError(t, err)

	w, err := r.Worktree()
	assert.NoError(t, err)

	filename := filepath.Join(dir, "foo")
	err = os.WriteFile(filename, []byte("bar!"), 0644)
	assert.NoError(t, err)

	co := &git.CommitOptions{
		All:               true,
		Author:            &object.Signature{Name: "test"},
		Committer:         &object.Signature{Name: "test"},
		AllowEmptyCommits: true,
	}

	_, err = w.Commit("foo\n", co)
	assert.NoError(t, err)

	// now run the create tag code
	tagOptions := Tag{
		Dir:              dir,
		PushTag:          false,
		FormattedVersion: "1.2.3",
	}

	err = tagOptions.TagRemote()
	assert.NoError(t, err)

	tags, err := r.TagObjects()
	assert.NoError(t, err)

	tag, err := tags.Next()
	assert.NoError(t, err)

	assert.Equal(t, "1.2.3", tag.Name)
}
