package command

import (
	"log/slog"

	"github.com/okt-limonikas/bublik-log-viewer/internal/constants"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print current version information",
	Long:  "Print current version of binary",
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("version information",
			"version", constants.Version,
			"date", constants.Date,
			"commit", constants.Commit)
	},
}
