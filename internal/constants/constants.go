package constants

import (
	"os/user"
	"path/filepath"
	"time"
)

const (
	BinName    = "blv"
	BuildPath  = "build"
	GithubRepo = "github.com/okt-limonikas/bublik-log-viewer"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

const UpdateInterval = time.Hour * 6

var LastUpdateFile string

func init() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	homeDir := currentUser.HomeDir

	LastUpdateFile = filepath.Join(homeDir, ".blv_last_update.txt")
}
