package search

import (
	"fmt"
	"strings"

	"capsule-package-manager/models"

	"charm.land/log/v2"
)

// Filters packages using SearchTerm and prints each matching package
// Called by main() in main.go; uses PackageSearchLoop() here
func BasicPackageSearch(SearchTerm string, BasicPackageInfo []models.PackageInfo) {
	searchResults := PackageSearchLoop(SearchTerm, BasicPackageInfo)
	for _, packageName := range searchResults {
		fmt.Printf("Name: %s\nVersion: %s\nRepo: %s\nDescription: %s\n",
			packageName.Name, packageName.Version, packageName.Repo, packageName.Desc)
	}
}

// Filters packages and prints detailed information for each match
// Called by main() in main.go; uses PackageSearchLoop() here
func DetailedPackageSearch(SearchTerm string, DetailedPackageInfo []models.PackageInfo) error {
	searchResults := PackageSearchLoop(SearchTerm, DetailedPackageInfo)
	for _, packageName := range searchResults {
		fmt.Printf("%-16s : %s\n", "From Repo", packageName.Repo)
		fmt.Printf("%-16s : %s\n", "Name", packageName.Name)
		fmt.Printf("%-16s : %s\n", "Version", packageName.Version)
		fmt.Printf("%-16s : %s\n", "Description", packageName.Desc)
		fmt.Printf("%-16s : %s\n", "Architecture", packageName.Architecture)
		fmt.Printf("%-16s : %s\n", "URL", packageName.URL)
		fmt.Printf("%-16s : %s\n", "Licenses", packageName.Licenses)
		fmt.Printf("%-16s : %s\n", "Groups", packageName.Groups)
		fmt.Printf("%-16s : %s\n", "Provides", packageName.Provides)
		fmt.Printf("%-16s : %s\n", "Depends On", packageName.DependsOn)
		fmt.Printf("%-16s : %s\n", "Optional Deps", packageName.OptionalDeps)
		fmt.Printf("%-16s : %s\n", "Conflicts With", packageName.ConflictsWith)
		fmt.Printf("%-16s : %s\n", "Replaces", packageName.Replaces)
		fmt.Printf("%-16s : %s\n", "Download Size", packageName.DownloadSize)
		fmt.Printf("%-16s : %s\n", "Installed Size", packageName.InstalledSize)
		fmt.Printf("%-16s : %s\n", "Packager", packageName.Packager)
		fmt.Printf("%-16s : %s\n", "Build Date", packageName.BuildDate)
		fmt.Printf("%-16s : %s\n", "Validated By", packageName.ValidatedBy)
		fmt.Println(" ")
	}
	return nil
}

// Filters packages by SearchTerm against name and description
// Gets SearchTerm and packages from both search functions in this file
// Returns searchResults to the calling search function
func PackageSearchLoop(SearchTerm string, packages []models.PackageInfo) []models.DetailedPackageInfo {
	log.Info("Searching for package", "package", SearchTerm)
	searchResults := []models.DetailedPackageInfo{}
	for _, packageName := range packages {
		if strings.Contains(packageName.Name, SearchTerm) || strings.Contains(packageName.Desc, SearchTerm) {
			searchResults = append(searchResults, packageName.DetailedPackageInfo)
		}
	}
	return searchResults
}
