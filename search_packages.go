package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"golang.org/x/term"
)

// Real search func but with placeholder code to test functionality before actually searching the repos
func BasicPackageSearch(SearchTerm string, BasicPackageInfo []DetailedPackageInfo) {
	searchResults := PackageSearchLoop(SearchTerm, BasicPackageInfo)
	tuiList := SearchList(searchResults)
	tui := tea.NewProgram(tuiList)
	tui.Run()
}

// Had to do it with seperate print statements as a multi-line print fucked with the formatting. Sadge
func DetailedPackageSearch(SearchTerm string, DetailedPackageInfo []DetailedPackageInfo) error {
	searchResults := PackageSearchLoop(SearchTerm, DetailedPackageInfo)

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return fmt.Errorf("could not get terminal width: %s", err)
	}

	detailView, err := DetailView(searchResults, 0, width, height)
	if err != nil {
		return fmt.Errorf("failed to get DetailView: %s", err)
	}
	detailView.standalone = true

	tea.NewProgram(detailView).Run()

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

func PackageSearchLoop(SearchTerm string, BasicPackageInfo []DetailedPackageInfo) []DetailedPackageInfo {
	searchResults := []DetailedPackageInfo{}
	for i, packageName := range BasicPackageInfo {
		if packageName.Name == SearchTerm {
			searchResults = append(searchResults, BasicPackageInfo[i])
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

// Saves including word wrap code in the print statements
func detailedWrappingHelper(packagesToWrap string, width int) (string, error) {
	var builder strings.Builder

	availableWidth := width - 19
	packageWords := strings.Fields(packagesToWrap)
	lineToWrap := ""

	for _, word := range packageWords {
		if len(lineToWrap)+len(word)+1 > availableWidth {
			builder.WriteString(lineToWrap)
			builder.WriteString("\n")
			builder.WriteString(strings.Repeat(" ", 19))
			lineToWrap = word + " "
			continue
		} else {
			lineToWrap += word + " "
		}
	}
	builder.WriteString(lineToWrap)
	builder.WriteString("\n")
	return builder.String(), nil
}

func showNoneWhereNoData(isNone string) string {
	if isNone == "" {
		return "None"
	}
	return isNone
}

func DetailedPackageOutput(SearchTerm string, DetailedPackageInfo []DetailedPackageInfo, width int) (string, error) {
	var builder strings.Builder

	for i, packageName := range DetailedPackageInfo {
		if packageName.Name == SearchTerm {
			wrappedDeps, err := detailedWrappingHelper(DetailedPackageInfo[i].DependsOn, width)
			if err != nil {
				return "", fmt.Errorf("line wrap helper function failed: %s", err)
			}
			fmt.Fprintf(&builder, "%-16s : %s\n", "From Repo", DetailedPackageInfo[i].Repo)
			fmt.Fprintf(&builder, "%-16s : %s\n", "Name", DetailedPackageInfo[i].Name)
			fmt.Fprintf(&builder, "%-16s : %s\n", "Version", DetailedPackageInfo[i].Version)
			fmt.Fprintf(&builder, "%-16s : %s\n", "Description", DetailedPackageInfo[i].Desc)
			fmt.Fprintf(&builder, "%-16s : %s\n", "Architecture", DetailedPackageInfo[i].Architecture)
			fmt.Fprintf(&builder, "%-16s : %s\n", "URL", DetailedPackageInfo[i].URL)
			fmt.Fprintf(&builder, "%-16s : %s\n", "Licenses", DetailedPackageInfo[i].Licenses)
			fmt.Fprintf(&builder, "%-16s : %s\n", "Groups", showNoneWhereNoData(DetailedPackageInfo[i].Groups))
			fmt.Fprintf(&builder, "%-16s : %s\n", "Provides", DetailedPackageInfo[i].Provides)
			//fmt.Fprintf(&builder, "%-16s : ", "Depends On")
			fmt.Fprintf(&builder, "%-16s : %s\n", "Depends On", wrappedDeps)
			fmt.Fprintf(&builder, "%-16s : %s\n", "Optional Deps", DetailedPackageInfo[i].OptionalDeps)
			fmt.Fprintf(&builder, "%-16s : %s\n", "Conflicts With", showNoneWhereNoData(DetailedPackageInfo[i].ConflictsWith))
			fmt.Fprintf(&builder, "%-16s : %s\n", "Replaces", showNoneWhereNoData(DetailedPackageInfo[i].Replaces))
			fmt.Fprintf(&builder, "%-16s : %s\n", "Download Size", DetailedPackageInfo[i].DownloadSize)
			fmt.Fprintf(&builder, "%-16s : %s\n", "Installed Size", DetailedPackageInfo[i].InstalledSize)
			fmt.Fprintf(&builder, "%-16s : %s\n", "Packager", DetailedPackageInfo[i].Packager)
			fmt.Fprintf(&builder, "%-16s : %s\n", "Build Date", DetailedPackageInfo[i].BuildDate)
			fmt.Fprintf(&builder, "%-16s : %s\n", "Validated By", DetailedPackageInfo[i].ValidatedBy)
			fmt.Fprintln(&builder)
		}
	}
	return builder.String(), nil
}
