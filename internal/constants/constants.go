package constants

import (
	"os"
	"os/user"
	"path/filepath"
	"time"

	"log/slog"
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
		slog.Error("failed to get current user", "error", err)
		os.Exit(1)
	}

	homeDir := currentUser.HomeDir
	LastUpdateFile = filepath.Join(homeDir, ".blv_last_update.txt")
}
