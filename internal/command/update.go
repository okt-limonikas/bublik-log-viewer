package command

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/Masterminds/semver"
	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
	"github.com/okt-limonikas/bublik-log-viewer/internal/constants"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update to latest version",
	Run: func(cmd *cobra.Command, args []string) {
		err := performUpdate()
		if err != nil {
			slog.Error("update failed", "error", err)
		}

		err = updateLastUpdateTime()
		if err != nil {
			slog.Error("failed to update last update file", "error", err)
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
		slog.Info("update check skipped - minimum update interval not reached")
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

	if shouldUpdate(versions.current, versions.latest) {
		fmt.Println("╭──────────────────────────────────────────")
		fmt.Println("│ 🆕 Update Available!")
		fmt.Println("│")
		fmt.Printf("│ 📦 Current version: v%s\n", versions.current)
		fmt.Printf("│ ⭐ Latest version: v%s\n", versions.latest)
		fmt.Println("│")
		fmt.Println("│ 💡 To update, run:")
		fmt.Println("│    blv self-update")
		fmt.Println("╰──────────────────────────────────────────")
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
	fmt.Println("╭──────────────────────────────────────────")
	fmt.Println("│ 🔄 Starting Update Process")
	fmt.Println("│")
	fmt.Printf("│ 📦 Current version: v%s\n", constants.Version)
	fmt.Printf("│ 📅 Build date: %s\n", constants.Date)
	fmt.Printf("│ 🔨 Commit: %s\n", constants.Commit)
	fmt.Println("│")
	fmt.Println("│ 🔍 Checking latest version...")

	latestVersion, err := updateManager.GetLatestVersion()
	if err != nil {
		return err
	}
	fmt.Printf("│ ✨ Latest version found: v%s\n", latestVersion)

	versions, err := getVersions()
	if err != nil {
		return fmt.Errorf("failed to create versions struct %w", err)
	}

	if !shouldUpdate(versions.current, versions.latest) {
		fmt.Println("│")
		fmt.Println("│ ✅ You're already on the latest version!")
		fmt.Println("╰──────────────────────────────────────────")
		return nil
	}

	fmt.Println("│")
	fmt.Printf("│ 📥 Downloading update (%s)...\n", archiveName)
	_, err = updateManager.Update()
	if err != nil {
		return err
	}

	fmt.Println("│")
	fmt.Printf("│ 🎉 Update successful! Now running v%s\n", latestVersion)
	fmt.Println("╰──────────────────────────────────────────")

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
