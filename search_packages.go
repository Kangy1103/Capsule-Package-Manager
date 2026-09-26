package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// Real search func but with placeholder code to test functionality before actually searching the repos
func BasicPackageSearch(SearchTerm string, BasicPackageInfo []DetailedPackageInfo) {
	for i, packageName := range BasicPackageInfo {
		if packageName.Name == SearchTerm {
			fmt.Printf("Name: %s\nVersion: %s\nRepo: %s\nDescription: %s\n",
				BasicPackageInfo[i].Name,
				BasicPackageInfo[i].Version,
				BasicPackageInfo[i].Repo,
				BasicPackageInfo[i].Desc,
			)
		}
	}
}

// Had to do it with seperate print statements as a multi-line print fucked with the formatting. Sadge
func DetailedPackageSearch(SearchTerm string, DetailedPackageInfo []DetailedPackageInfo) {
	for i, packageName := range DetailedPackageInfo {
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
	}
}

// Saves including word wrap code in the print statements
func detailedWrappingHelper(packagesToWrap string) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Could not get terminal width")
		return
	}
	availableWidth := width - 19
	packageWords := strings.Fields(packagesToWrap)
	lineToWrap := ""

	for _, word := range packageWords {
		if len(lineToWrap)+len(word) > availableWidth {
			fmt.Printf("%s\n", lineToWrap)
			fmt.Print(strings.Repeat(" ", 19))
			lineToWrap = word + " "
			continue
		} else {
			lineToWrap += word + " "
		}
	}
	fmt.Printf("%s\n", lineToWrap)
}

func showNoneWhereNoData(isNone string) string {
	if isNone == "" {
		return "None"
	}
	return isNone
}
