package tui

import (
	"fmt"
	"strings"

	"capsule-package-manager/models"
)

// Saves including word wrap code in the print statements
func DetailedWrappingHelper(packagesToWrap string, width int) (string, error) {
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

func ShowNoneWhereNoData(isNone string) string {
	if isNone == "" {
		return "None"
	}
	return isNone
}

func DetailedPackageOutput(SearchTerm string, DetailedPackageInfo []models.DetailedPackageInfo, width int) (string, error) {
	var builder strings.Builder

	for i, packageName := range DetailedPackageInfo {
		if packageName.Name == SearchTerm {
			wrappedDeps, err := DetailedWrappingHelper(DetailedPackageInfo[i].DependsOn, width)
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
			fmt.Fprintf(&builder, "%-16s : %s\n", "Groups", ShowNoneWhereNoData(DetailedPackageInfo[i].Groups))
			fmt.Fprintf(&builder, "%-16s : %s\n", "Provides", DetailedPackageInfo[i].Provides)
			//fmt.Fprintf(&builder, "%-16s : ", "Depends On")
			fmt.Fprintf(&builder, "%-16s : %s\n", "Depends On", wrappedDeps)
			fmt.Fprintf(&builder, "%-16s : %s\n", "Optional Deps", DetailedPackageInfo[i].OptionalDeps)
			fmt.Fprintf(&builder, "%-16s : %s\n", "Conflicts With", ShowNoneWhereNoData(DetailedPackageInfo[i].ConflictsWith))
			fmt.Fprintf(&builder, "%-16s : %s\n", "Replaces", ShowNoneWhereNoData(DetailedPackageInfo[i].Replaces))
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
