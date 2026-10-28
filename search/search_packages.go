package search

import (
	"fmt"
	"os"
	"strings"

	"capsule-package-manager/models"
	"capsule-package-manager/tui"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"golang.org/x/term"
)

// Real search func but with placeholder code to test functionality before actually searching the repos
func BasicPackageSearch(SearchTerm string, BasicPackageInfo []models.PackageInfo) {
	searchResults := PackageSearchLoop(SearchTerm, BasicPackageInfo)
	tuiList := tui.SearchList(searchResults)
	tui := tea.NewProgram(tuiList)
	log.Debug("Opening search list")
	if _, err := tui.Run(); err != nil {
		log.Error("could not open TUI", "tui", err)
	}
}

// Had to do it with seperate print statements as a multi-line print fucked with the formatting. Sadge
func DetailedPackageSearch(SearchTerm string, DetailedPackageInfo []models.PackageInfo) error {
	searchResults := PackageSearchLoop(SearchTerm, DetailedPackageInfo)

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return fmt.Errorf("could not get terminal width: %s", err)
	}

	detailView, err := tui.DetailView(searchResults, 0, width, height)
	if err != nil {
		return fmt.Errorf("failed to get DetailView: %s", err)
	}
	detailView.Standalone = true

	log.Debug("Opening detailed package view")
	if _, err := tea.NewProgram(detailView).Run(); err != nil {
		log.Error("could not open detailed search view", "tui", err)
	}

	/*for i, packageName := range DetailedPackageInfo {
		if packageName.Name == SearchTerm {
			fmt.Printf("%-16s : %s\n", "From Repo", DetailedPackageInfo[i].Repo)
			fmt.Printf("%-16s : %s\n", "Name", DetailedPackageInfo[i].Name)
			fmt.Printf("%-16s : %s\n", "Version", DetailedPackageInfo[i].Version)
			fmt.Printf("%-16s : %s\n", "Description", DetailedPackageInfo[i].Desc)
			fmt.Printf("%-16s : %s\n", "Architecture", DetailedPackageInfo[i].Architecture)
			fmt.Printf("%-16s : %s\n", "URL", DetailedPackageInfo[i].URL)
			fmt.Printf("%-16s : %s\n", "Licenses", DetailedPackageInfo[i].Licenses)
			fmt.Printf("%-16s : %s\n", "Groups", showNoneWhereNoData(DetailedPackageInfo[i].Groups))
			fmt.Printf("%-16s : %s\n", "Provides", DetailedPackageInfo[i].Provides)
			fmt.Printf("%-16s : ", "Depends On")
			detailedWrappingHelper(DetailedPackageInfo[i].DependsOn)
			fmt.Printf("%-16s : %s\n", "Optional Deps", DetailedPackageInfo[i].OptionalDeps)
			fmt.Printf("%-16s : %s\n", "Conflicts With", showNoneWhereNoData(DetailedPackageInfo[i].ConflictsWith))
			fmt.Printf("%-16s : %s\n", "Replaces", showNoneWhereNoData(DetailedPackageInfo[i].Replaces))
			fmt.Printf("%-16s : %s\n", "Download Size", DetailedPackageInfo[i].DownloadSize)
			fmt.Printf("%-16s : %s\n", "Installed Size", DetailedPackageInfo[i].InstalledSize)
			fmt.Printf("%-16s : %s\n", "Packager", DetailedPackageInfo[i].Packager)
			fmt.Printf("%-16s : %s\n", "Build Date", DetailedPackageInfo[i].BuildDate)
			fmt.Printf("%-16s : %s\n", "Validated By", DetailedPackageInfo[i].ValidatedBy)
			fmt.Println(" ")
		}
	}*/
	return nil
}

func PackageSearchLoop(SearchTerm string, packages []models.PackageInfo) []models.DetailedPackageInfo {
	log.Info("Searching for package", "package", SearchTerm)
	searchResults := []models.DetailedPackageInfo{}
	for _, packageName := range packages {
		if strings.Contains(packageName.Name, SearchTerm) || strings.Contains(packageName.Desc, SearchTerm) {
			searchResults = append(searchResults, packageName.DetailedPackageInfo)
			/*fmt.Printf("Name: %s\nVersion: %s\nRepo: %s\nDescription: %s\n",
				BasicPackageInfo[i].Name,
				BasicPackageInfo[i].Version,
				BasicPackageInfo[i].Repo,
				BasicPackageInfo[i].Desc,
			)*/
		}
	}
	return searchResults
}
