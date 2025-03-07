package command

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/okt-limonikas/bublik-log-viewer/internal/constants"
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
	// Set up structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	err := rootCmd.Execute()
	if err != nil {
		slog.Error("command execution failed", "error", err)
		os.Exit(1)
	}
}
