package cmd

import (
	"fmt"
	"log"

	"github.com/okt-limonikas/bublik-log-viewer/constants"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   constants.BinName,
	Short: "Log viewer for JSON logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && cmd.Flags().NFlag() == 0 {
			return fmt.Errorf("no command specified")
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
