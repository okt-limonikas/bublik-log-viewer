package cmd

import (
	"fmt"
	"log"
	"runtime"
	"time"

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

	log.Println("Determining if can update...")
	canUpdate, err := updater.CanUpdate()
	if err != nil {
		return nil, err
	}

	if !canUpdate {
		log.Println("Can't update")
		log.Println("Skipping...")
		return nil, nil
	}

	log.Printf("Loading \"%s\"", archiveName)
	if _, err = updater.Update(); err != nil {
		return nil, err
	}

	log.Printf("Succesfully updated to %s", latestVersion)

	return nil, nil
}

type ReleaseResponse struct {
	URL             string    `json:"url"`
	AssetsURL       string    `json:"assets_url"`
	UploadURL       string    `json:"upload_url"`
	HTMLURL         string    `json:"html_url"`
	ID              int64     `json:"id"`
	Author          Author    `json:"author"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Assets          []Asset   `json:"assets"`
	TarballURL      string    `json:"tarball_url"`
	ZipballURL      string    `json:"zipball_url"`
	Body            string    `json:"body"`
}

type Asset struct {
	URL                string    `json:"url"`
	ID                 int64     `json:"id"`
	NodeID             string    `json:"node_id"`
	Name               string    `json:"name"`
	Label              string    `json:"label"`
	Uploader           Author    `json:"uploader"`
	ContentType        string    `json:"content_type"`
	State              string    `json:"state"`
	Size               int64     `json:"size"`
	DownloadCount      int64     `json:"download_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	BrowserDownloadURL string    `json:"browser_download_url"`
}

type Author struct {
	Login             string `json:"login"`
	ID                int64  `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}
