package main

import (
	"fmt"
	"os"
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

// Switch statement needs optimising when more commands are added
// For exmaple, SearchTerm coming from elsewhere instead of being read for every case
// Will also need an error check for os.Args[1]. Possibly showing a help message if "capsule" is ran on it's own
func main() {
	PackageInfo, err := FetchPackages()
	if err != nil {
		fmt.Println(err)
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
