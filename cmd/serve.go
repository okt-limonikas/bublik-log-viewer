package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/okt-limonikas/bublik-log-viewer/constants"
	"github.com/okt-limonikas/bublik-log-viewer/frontend"
	"github.com/okt-limonikas/bublik-log-viewer/utils"
	"github.com/spf13/cobra"
)

// Host for created HTTP server
var host string

// Port for create HTTP server
var port string

// Determines if should try to open user browser
var shouldOpenBrowser bool

// Determines if passed value is URL
var isRemote bool

var serveLogsCmd = &cobra.Command{
	Use:   "serve <log_path_or_url>",
	Short: "Serve JSON logs from a specified directory or URL",
	Long:  `This command will serve logs from the specified directory or URL`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var path string = args[0]

		input := LogsInput{path, host, port, shouldOpenBrowser, isRemote}

		err := checkForUpdateScheduled()
		if err != nil {
			log.Println(err)
		}

		serveLogs(&input)
	},
}

func init() {
	rootCmd.AddCommand(serveLogsCmd)

	serveLogsCmd.Flags().StringVar(&host, "host", "127.0.0.1", "Host address")
	serveLogsCmd.Flags().StringVar(&port, "port", "5050", "Port number")
	serveLogsCmd.Flags().BoolVar(&shouldOpenBrowser, "open", false, "Should open browser")
	serveLogsCmd.Flags().BoolVar(&isRemote, "remote", false, "Serve remote logs")
}

type LogsInput struct {
	path              string
	host              string
	port              string
	shouldOpenBrowser bool
	isRemote          bool
}

func createLogHandler(pathOrUrl string, isRemote bool) http.Handler {
	var handleJsonLogs http.Handler

	if isRemote {
		url, err := url.Parse(pathOrUrl)
		if err != nil {
			log.Fatal("Failed to parse URL")
		}

		if !strings.HasPrefix(url.Scheme, "http://") && !strings.HasPrefix(url.Scheme, "https://") {
			log.Fatal("Not supported URL scheme")
		}

		handleJsonLogs = httputil.NewSingleHostReverseProxy(url)
	} else {
		fs := http.FileServer(http.Dir(pathOrUrl))
		handleJsonLogs = http.StripPrefix("/json/", fs)
	}

	return handleJsonLogs
}

func resolvePath(pathOrUrl string, isRemote bool) string {
	var resolvedPath string

	if isRemote {
		resolvedPath = pathOrUrl
	} else {
		path, err := filepath.Abs(pathOrUrl)
		if err != nil {
			log.Fatal("Error getting absolute path of log directory:", err)
		}

		resolvedPath = path
	}

	return resolvedPath
}

func maybeOpenBrowser(url string, shouldOpenBrowser bool) {
	if !shouldOpenBrowser {
		return
	}

	err := utils.OpenURL(url)
	if err != nil {
		log.Printf("Failed to open browser!")
	}
}

func serveLogs(input *LogsInput) {
	var resolvedUrl = fmt.Sprintf("http://%s:%s", input.host, input.port)

	http.Handle("/json/", createLogHandler(input.path, input.isRemote))
	http.HandleFunc("/", handleSPA)

	log.Printf("The server is listening at %s", resolvedUrl)
	log.Printf("Serving log files from %s", resolvePath(input.path, input.isRemote))

	maybeOpenBrowser(resolvedUrl, input.shouldOpenBrowser)

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", input.host, input.port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleSPA(w http.ResponseWriter, r *http.Request) {
	f, err := frontend.BuildFs.Open(filepath.Join(constants.BUILD_PATH, r.URL.Path))

	if os.IsNotExist(err) {
		index, err := frontend.BuildFs.ReadFile(filepath.Join(constants.BUILD_PATH, "index.html"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write(index)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	http.FileServer(frontend.BuildHTTPFS()).ServeHTTP(w, r)
}
