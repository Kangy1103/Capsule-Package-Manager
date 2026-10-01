package main

import (
	"fmt"
	"os"

	"charm.land/log/v2"
)

type BasicPackageInfo struct {
	Name    string
	Version string
	Desc    string
	Repo    string
}

type DetailedPackageInfo struct {
	BasicPackageInfo
	Architecture  string
	URL           string
	Licenses      string
	Groups        string
	Provides      string
	DependsOn     string
	OptionalDeps  string
	ConflictsWith string
	Replaces      string
	DownloadSize  string
	InstalledSize string
	Packager      string
	BuildDate     string
	ValidatedBy   string
}

type Options struct {
	ReportCaller    bool
	ReportTimestamp bool
}

const DefaultTimeFormat = "2006/01/02 15:04:05"

// Switch statement needs optimising when more commands are added
// For exmaple, SearchTerm coming from elsewhere instead of being read for every case
// Will also need an error check for os.Args[1]. Possibly showing a help message if "capsule" is ran on it's own
func main() {
	logger, closeLog, err := Logging()
	if err != nil {
		fmt.Println("unable to start logger")
		return
	}
	defer closeLog()
	log.SetDefault(logger)

	log.Info("Starting Capsule Package Manager")

	PackageInfo, err := FetchPackages()
	if err != nil {
		log.Errorf("unable to fetch packages: %v", err)
		return
	}

	switch os.Args[1] {
	case "search":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a package name to search")
			return
		}
		SearchTerm := os.Args[2]
		BasicPackageSearch(SearchTerm, PackageInfo)

	case "details":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a package name to search")
			return
		}
		SearchTerm := os.Args[2]
		DetailedPackageSearch(SearchTerm, PackageInfo)
	}
}

type closeFunc func() error

func Logging() (*log.Logger, closeFunc, error) {
	logFile := "logs/capsule.log"
	err := os.MkdirAll("logs", 0o755)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create log folder: %v", err)
	}
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open log file: %v", err)
	}
	/*bufferedFile := bufio.NewWriterSize(file, 8192)
	close := func() error {
		if err := bufferedFile.Flush(); err != nil {
			return fmt.Errorf("could not flush file: %v", err)
		}
		if err := file.Close(); err != nil {
			return fmt.Errorf("could not close file: %v", err)
		}

		return nil
	}*/
	close := func() error {
		if err := file.Close(); err != nil {
			return fmt.Errorf("could not close file: %v", err)
		}

		return nil
	}

	logger := log.NewWithOptions(file, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		Level:           log.DebugLevel,
	})
	return logger, close, nil
}
