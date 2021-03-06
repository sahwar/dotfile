package dotfile

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUncompressRevision(t *testing.T) {
	t.Run("uncompress revision error", func(t *testing.T) {
		s := &MockStorer{revisionErr: true}

		_, err := UncompressRevision(s, testHash)
		assert.Error(t, err)
	})

	t.Run("uncompress error", func(t *testing.T) {
		s := &MockStorer{uncompressErr: true}
		_, err := UncompressRevision(s, testHash)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		_, err := UncompressRevision(new(MockStorer), testHash)
		assert.NoError(t, err)
	})

}

func TestIsClean(t *testing.T) {
	t.Run("content error", func(t *testing.T) {
		s := &MockStorer{dirtyContentErr: true}
		_, err := IsClean(s, testHash)
		assert.Error(t, err)
	})
	t.Run("true with no dirty content", func(t *testing.T) {
		s := &MockStorer{noDirtyContent: true}
		clean, err := IsClean(s, testHash)
		assert.NoError(t, err)
		assert.True(t, clean)
	})
	t.Run("false", func(t *testing.T) {
		s := &MockStorer{}
		clean, err := IsClean(s, testHash)
		assert.NoError(t, err)
		assert.False(t, clean)
	})
}

func TestDiff(t *testing.T) {
	t.Run("uncompress error", func(t *testing.T) {
		s := &MockStorer{uncompressErr: true}
		_, err := Diff(s, testHash, testHash)
		assert.Error(t, err)
	})

	t.Run("get content error", func(t *testing.T) {
		s := &MockStorer{dirtyContentErr: true}
		_, err := Diff(s, testHash, "")
		assert.Error(t, err)
	})

	t.Run("no changes error", func(t *testing.T) {
		s := &MockStorer{}
		_, err := Diff(s, testHash, testHash)
		assert.True(t, errors.Is(err, ErrNoChanges))
	})

}
