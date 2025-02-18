package command

import (
	"fmt"

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
		fmt.Println("╭──────────────────────────────────────────")
		fmt.Println("│ 📦 Bublik Log Viewer")
		fmt.Println("│")
		fmt.Printf("│ 🔖 Version: %s\n", constants.Version)
		fmt.Printf("│ 📅 Build date: %s\n", constants.Date)
		fmt.Printf("│ 🔨 Commit: %s\n", constants.Commit)
		fmt.Println("╰──────────────────────────────────────────")
	},
}
