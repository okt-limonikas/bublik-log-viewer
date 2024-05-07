package constants

import (
	"os/user"
	"path/filepath"
	"time"
)

const BIN_NAME = "bublik-log-viewer"
const BUILD_PATH = "build"
const GITHUB_REPO = "github.com/okt-limonikas/bublik-log-viewer"

var Version = "dev"
var Commit = "none"
var Date = "unknown"

const UpdateInterval = time.Hour * 6

var LastUpdateFile string

func init() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	homeDir := currentUser.HomeDir

	LastUpdateFile = filepath.Join(homeDir, ".last_update.txt")
}
