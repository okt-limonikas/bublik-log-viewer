package cmd

import (
	"log"

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
		log.Printf("Current version: %s", version)
		log.Printf("Build date: %s", date)
		log.Printf("Commit: %s", commit)
	},
}
