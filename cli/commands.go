package cli

import (
	"github.com/knoebber/dotfile/file"
	"gopkg.in/alecthomas/kingpin.v2"

	"fmt"
	"os"
)

const (
	defaultStorageDir  string = ".dotfile/"
	defaultStorageName string = "files.json"
)

// Dotfile depends on the system having the concept of a home directory.
func getHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home
}

func AddCommandsToApplication(app *kingpin.Application) {
	storage := &file.Storage{
		Home: getHome(),
	}

	app.Flag("storage-dir", "The directory where version control storage is stored").
		Default(fmt.Sprintf("%s/%s", storage.Home, defaultStorageDir)).
		StringVar(&storage.Dir)
	app.Flag("storage-name", "The main json file that tracks checked in files").
		Default(defaultStorageName).
		StringVar(&storage.Name)

	addInitSubCommandToApplication(app, storage)
	addEditSubCommandToApplication(app, storage)
	addDiffSubCommandToApplication(app, storage)
	addLogSubCommandToApplication(app, storage)
	addCheckoutSubCommandToApplication(app, storage)
	addCommitSubCommandToApplication(app, storage)
	addPushSubCommandToApplication(app, storage)
	addPullSubCommandToApplication(app, storage)
}
