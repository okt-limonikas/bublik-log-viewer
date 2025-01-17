package cmd

import (
	"github.com/okt-limonikas/bublik-log-viewer/constants"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print current version information",
	Long:  "Print current version of binary",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Version: %s", constants.Version)
		log.Printf("Date: %s", constants.Date)
		log.Printf("Commit: %s", constants.Commit)
	},
}
