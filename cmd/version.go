package cmd

import (
	"log"

	"github.com/okt-limonikas/bublik-log-viewer/constants"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  "Print current version of binary",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Version: %s", constants.Version)
		log.Printf("Date: %s", constants.Date)
		log.Printf("Commit: %s", constants.Commit)
	},
}
