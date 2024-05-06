package cmd

import (
	"fmt"
	"log"
	"runtime"

	"github.com/Masterminds/semver"
	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
	"github.com/okt-limonikas/bublik-log-viewer/constants"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update to latest version",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := Update()
		if err != nil {
			log.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func Update() (interface{}, error) {
	log.Println("Starting update process...")

	log.Printf("Version: %s", constants.Version)
	log.Printf("Commit: %s", constants.Commit)
	log.Printf("Date: %s", constants.Date)

	archiveName := getFormattedBinaryInfo()

	updater := &updater.Updater{
		Provider: &provider.Github{
			RepositoryURL: constants.GITHUB_REPO,
			ArchiveName:   archiveName,
		},
		ExecutableName: constants.BIN_NAME,
		Version:        constants.Version,
	}

	log.Println("Checking latest version...")
	latestVersion, err := updater.GetLatestVersion()
	if err != nil {
		return nil, err
	}
	log.Printf("Latest version: %s", latestVersion)

	currentVersionSem, currentErr := semver.NewVersion(constants.Version)
	if currentErr != nil {
		log.Println("Failed to construct current version struct")
	}

	latestVersionSem, latestErr := semver.NewVersion(latestVersion)
	if latestErr != nil {
		log.Println("Failed to construct latest version struct")
	}

	if currentErr != nil || latestErr != nil {
		log.Println("Failed to construct semver. Trying to update...")
	}

	if currentErr == nil && latestErr == nil {
		if currentVersionSem.GreaterThan(latestVersionSem) || currentVersionSem.Equal(latestVersionSem) {
			log.Println("You are already on latest version...")
			return nil, nil
		}
	}

	log.Printf("Loading \"%s\"", archiveName)
	_, err = updater.Update()
	if err != nil {
		return nil, err
	}

	log.Printf("Succesfully updated to %s", latestVersion)

	return nil, nil
}

func getFormattedBinaryInfo() string {
	switch runtime.GOOS {
	case "windows":
		return fmt.Sprintf("%s_%s_%s.zip", constants.BIN_NAME, "Windows", runtime.GOARCH)
	case "darwin":
		return fmt.Sprintf("%s_%s_%s.tar.gz", constants.BIN_NAME, "Darwin", runtime.GOARCH)
	default:
		return fmt.Sprintf("%s_%s_%s.tar.gz", constants.BIN_NAME, "Linux", runtime.GOARCH)
	}
}
