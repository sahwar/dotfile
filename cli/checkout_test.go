package cli

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckout(t *testing.T) {
	clearTestStorage(t)
	initTestFile(t)

	checkoutCommand := new(checkoutCommand)

	t.Run("returns error when file is not tracked", func(t *testing.T) {
		checkoutCommand.alias = notTrackedFile
		assert.Error(t, checkoutCommand.run(nil))
	})

	t.Run("ok", func(t *testing.T) {
		checkoutCommand.alias = trackedFileAlias
		assert.NoError(t, checkoutCommand.run(nil))
	})

	t.Run("returns error when file has untracked changes", func(t *testing.T) {
		updateTestFile(t)
		checkoutCommand.alias = trackedFileAlias
		assert.Error(t, checkoutCommand.run(nil))
	})

	t.Run("ok to checkout deleted file", func(t *testing.T) {
		_ = os.Remove(trackedFile)
		checkoutCommand.alias = trackedFileAlias
		assert.NoError(t, checkoutCommand.run(nil))
	})

}
