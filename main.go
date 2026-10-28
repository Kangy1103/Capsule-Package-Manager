package main

import (
	"fmt"
	"os"

	packagemanagement "capsule-package-manager/package_management"
	"capsule-package-manager/search"

	"charm.land/log/v2"
)

// Switch statement needs optimising when more commands are added
// For exmaple, SearchTerm coming from elsewhere instead of being read for every case
// Will also need an error check for os.Args[1]. Possibly showing a help message if "capsule" is ran on it's own
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: capsule <command> [arguments]")
		return
	}
	logger, closeLog, err := Logging()
	if err != nil {
		fmt.Printf("unable to start logger: %v\n", err)
		return
	}
	defer closeLog()
	log.SetDefault(logger)

	log.Info("Starting Capsule Package Manager")

	cmd := os.Args[1]

	switch cmd {
	case "version", "help":
		if cmd == "version" {
			fmt.Println("Capsule v0.0.4")
		}

	case "search", "details", "install", "uninstall":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a package name")
			return
		}
		SearchTerm := os.Args[2]

		PackageInfo, handle, err := search.FetchPackages()
		if err != nil {
			log.Error("unable to fetch packages", "search", err)
			return
		}
		switch cmd {
		case "search":
			search.BasicPackageSearch(SearchTerm, PackageInfo)
		case "details":
			search.DetailedPackageSearch(SearchTerm, PackageInfo)
		case "install":
			packagemanagement.InstallPackage(SearchTerm, handle, PackageInfo)
		case "uninstall":
			packagemanagement.UninstallPackage()
		}
	default:
		fmt.Printf("unknown command: %v\n", cmd)
	}
}

type closeFunc func() error

func Logging() (*log.Logger, closeFunc, error) {
	logFolder := "/mnt/Goku/capsule/logs"
	logFile := "/mnt/Goku/capsule/logs/capsule.log"
	err := os.MkdirAll(logFolder, 0o755)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create log folder: %v", err)
	}
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open log file: %v", err)
	}

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
