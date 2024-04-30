package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ts-factory/bublik-log-viewer/frontend"
)

func main() {
	var fLogPath string
	var fPort string
	var fHost string

	flag.StringVar(&fHost, "host", "127.0.0.1", "Host address")
	flag.StringVar(&fPort, "port", "5050", "Port number")
	flag.StringVar(&fLogPath, "log-path", "", "Path to serve log files from")
	flag.Parse()

	if fLogPath == "" {
		fmt.Println("You have not provided --log-path flag so we serve files from default location: \"./json\"")
		flag.PrintDefaults()
	}

	var logPath string
	if fLogPath == "" {
		logPath = "./json"
	} else {
		logPath = fLogPath
	}

	fs := http.FileServer(http.Dir(logPath))
	http.Handle("/json/", http.StripPrefix("/json/", fs))
	http.HandleFunc("/", handleSPA)

	resolvedLogPath, err := filepath.Abs(logPath)
	if err != nil {
		fmt.Println("Error getting absolute path of log directory:", err)
		os.Exit(1)
	}

	log.Printf("The server is listening at http://%s:%s", fHost, fPort)
	log.Printf("Serving log files from %s", resolvedLogPath)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", fHost, fPort), nil))
}

func handleSPA(w http.ResponseWriter, r *http.Request) {
	buildPath := "build"
	f, err := frontend.BuildFs.Open(filepath.Join(buildPath, r.URL.Path))

	if os.IsNotExist(err) {
		index, err := frontend.BuildFs.ReadFile(filepath.Join(buildPath, "index.html"))

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
