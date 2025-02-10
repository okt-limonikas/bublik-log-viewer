package command

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/okt-limonikas/bublik-log-viewer/frontend"
	"github.com/okt-limonikas/bublik-log-viewer/internal/constants"
	"github.com/okt-limonikas/bublik-log-viewer/internal/utils"
	"github.com/spf13/cobra"
)

// Host for created HTTP server
var host string

// Port for create HTTP server
var port string

// Determines if you should try to open user browser
var shouldOpenBrowser bool

// Determines if passed value is URL
var isRemote bool

var serveLogsCmd = &cobra.Command{
	Use:   "serve <log_path_or_url>",
	Short: "Serve JSON logs from a specified directory or URL",
	Long:  `This command will serve logs from the specified directory or URL`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		input := LogsInput{path, host, port, shouldOpenBrowser, isRemote}

		err := checkForUpdateScheduled()
		if err != nil {
			slog.Error("error checking for update scheduled", "error", err)
		}

		serveLogs(&input)
	},
}

func init() {
	rootCmd.AddCommand(serveLogsCmd)

	serveLogsCmd.Flags().StringVar(&host, "host", "127.0.0.1", "Host address")
	serveLogsCmd.Flags().StringVar(&port, "port", "5050", "Port number")
	serveLogsCmd.Flags().BoolVar(&shouldOpenBrowser, "open", true, "Should open browser")
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
		parsedUrl, err := url.Parse(pathOrUrl)
		if err != nil {
			slog.Error("failed to parse URL", "error", err)
			os.Exit(1)
		}

		if !strings.HasPrefix(parsedUrl.Scheme, "http://") && !strings.HasPrefix(parsedUrl.Scheme, "https://") {
			slog.Error("not supported URL scheme")
			os.Exit(1)
		}

		handleJsonLogs = httputil.NewSingleHostReverseProxy(parsedUrl)
	} else {
		fileServer := http.FileServer(http.Dir(pathOrUrl))
		handleJsonLogs = http.StripPrefix("/json/", fileServer)
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
			slog.Error("error getting absolute path of log directory", "error", err)
			os.Exit(1)
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
		slog.Warn("failed to open browser", "error", err)
	}
}

func serveLogs(input *LogsInput) {
	resolvedUrl := fmt.Sprintf("http://%s:%s", input.host, input.port)

	http.Handle("/json/", createLogHandler(input.path, input.isRemote))
	http.HandleFunc("/", handleSPA)

	slog.Info("server listening", "url", resolvedUrl)
	slog.Info("serving log files", "path", resolvePath(input.path, input.isRemote))

	maybeOpenBrowser(resolvedUrl, input.shouldOpenBrowser)

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", input.host, input.port), nil)
	if err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}

func handleSPA(w http.ResponseWriter, r *http.Request) {
	f, err := frontend.BuildFs.Open(filepath.Join(constants.BuildPath, r.URL.Path))

	if os.IsNotExist(err) {
		index, err := frontend.BuildFs.ReadFile(filepath.Join(constants.BuildPath, "index.html"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		_, err = w.Write(index)
		if err != nil {
			return
		}
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(f fs.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)

	http.FileServer(frontend.BuildHTTPFS()).ServeHTTP(w, r)
}
