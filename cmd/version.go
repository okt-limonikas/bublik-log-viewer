package cmd

import (
	"fmt"

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
		version, err := getCurrentVersion()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%s", version)
	},
}

func getCurrentVersion() (string, error) {
	return "Helo", nil
}
