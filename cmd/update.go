package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/Masterminds/semver"
	"github.com/charmbracelet/lipgloss"
	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
	"github.com/okt-limonikas/bublik-log-viewer/constants"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update to latest version",
	Run: func(cmd *cobra.Command, args []string) {
		err := performUpdate()
		if err != nil {
			log.Fatal(err)
		}

		err = updateLastUpdateTime()
		if err != nil {
			log.Fatal("Failed to update last update file")
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

var archiveName = getArchiveName()

var updateManager = &updater.Updater{
	Provider: &provider.Github{
		RepositoryURL: constants.GithubRepo,
		ArchiveName:   archiveName,
	},
	ExecutableName: constants.BinName,
	Version:        constants.Version,
}

func readLastUpdateTime() (time.Time, error) {
	if _, err := os.Stat(constants.LastUpdateFile); os.IsNotExist(err) {
		// Create the file if it doesn't exist
		err := os.WriteFile(constants.LastUpdateFile, []byte(time.Now().Format(time.RFC3339)), 0644)
		if err != nil {
			return time.Now(), fmt.Errorf("failed to create last update time file: %w", err)
		}
		return time.Now(), nil
	}

	data, err := os.ReadFile(constants.LastUpdateFile)
	if err != nil {
		return time.Now(), fmt.Errorf("failed to read last update time from file: %w", err)
	}

	lastUpdate, err := time.Parse(time.RFC3339, string(data))
	if err != nil {
		return time.Now(), fmt.Errorf("failed to parse last update time: %w", err)
	}

	return lastUpdate, nil
}

func checkForUpdateScheduled() error {
	lastUpdate, err := readLastUpdateTime()
	if err != nil {
		return fmt.Errorf("failed to get last update time: %w", err)
	}

	// If the last update check was less than UpdateInterval ago, do not check again
	if time.Since(lastUpdate) < constants.UpdateInterval {
		log.Println("Update check skipped. Minimum update interval not reached.")
		return nil
	}

	versions, err := getVersions()
	if err != nil {
		return fmt.Errorf("failed to get versions %w", err)
	}

	err = updateLastUpdateTime()
	if err != nil {
		return fmt.Errorf("failed to update last update time: %w", err)
	}

	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}).
		Padding(1).Border(lipgloss.NormalBorder())

	if shouldUpdate(versions.current, versions.latest) {
		message := fmt.Sprintf("New version available: v%s\nYou can update: `bublik-log-viewer update`", versions.latest.String())
		fmt.Println(style.Render(message))
	}

	return nil
}

func updateLastUpdateTime() error {
	err := os.WriteFile(constants.LastUpdateFile, []byte(time.Now().Format(time.RFC3339)), 0644)
	if err != nil {
		return fmt.Errorf("failed to update last update time file: %w", err)
	}
	return nil
}

type Versions struct {
	current *semver.Version
	latest  *semver.Version
}

func getVersions() (*Versions, error) {
	latestVersion, err := updateManager.GetLatestVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest version: %w", err)
	}

	latestVersionSem, err := semver.NewVersion(latestVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to create latest version struct: %w", err)
	}

	currentVersion, err := semver.NewVersion(constants.Version)
	if err != nil {
		currentVersion, _ = semver.NewVersion("0.0.0")
	}

	return &Versions{current: currentVersion, latest: latestVersionSem}, nil
}

func shouldUpdate(current *semver.Version, latest *semver.Version) bool {
	return latest.GreaterThan(current)
}

func performUpdate() error {
	log.Println("Starting update process...")

	log.Printf("Version: %s", constants.Version)
	log.Printf("Commit: %s", constants.Commit)
	log.Printf("Date: %s", constants.Date)

	log.Println("Checking latest version...")
	latestVersion, err := updateManager.GetLatestVersion()
	if err != nil {
		return err
	}
	log.Printf("Latest version: %s", latestVersion)

	versions, err := getVersions()
	if err != nil {
		return fmt.Errorf("failed to create versions struct %w", err)
	}

	if !shouldUpdate(versions.current, versions.latest) {
		log.Println("You are already on latest version...")
		return nil
	}

	log.Printf("Loading \"%s\"", archiveName)
	_, err = updateManager.Update()
	if err != nil {
		return err
	}

	log.Printf("Succesfully updated to %s", latestVersion)

	return nil
}

func getArchiveName() string {
	switch runtime.GOOS {
	case "windows":
		return fmt.Sprintf("%s_%s_%s.zip", constants.BinName, "Windows", runtime.GOARCH)
	case "darwin":
		return fmt.Sprintf("%s_%s_%s.tar.gz", constants.BinName, "Darwin", runtime.GOARCH)
	default:
		return fmt.Sprintf("%s_%s_%s.tar.gz", constants.BinName, "Linux", runtime.GOARCH)
	}
}
