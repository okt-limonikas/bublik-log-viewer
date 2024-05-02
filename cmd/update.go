package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update to latest version",
	Run: func(cmd *cobra.Command, args []string) {
		Update()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func Update() {
	// TODO:
	log.Println("Starting update process...")
}
