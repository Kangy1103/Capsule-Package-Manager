package tui

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type InstallSummaryRow struct {
	Repo             string
	Name             string
	InstalledVersion string
	NewVersion       string
	DownloadSize     int64
	InstalledSize    int64
	PreviousSize     int64
}

type InstallSummary struct {
	Add    []InstallSummaryRow
	Remove []InstallSummaryRow
}

func RenderInstallSummary(summary InstallSummary) string {
	lines := []string{}
	downloadSize, installedSize, previousSize := installSummaryTotals(summary)
	downloadString := fmt.Sprintf("Total download: %v", mib(downloadSize))
	installedString := fmt.Sprintf("Installed size: %v", mib(installedSize))
	previousString := fmt.Sprintf("Net disk change: %v", signed(installedSize-previousSize))
	installs, upgrades, reinstalls, packages := installedSummaryCounts(summary)
	lines = append(lines, "Capsule / Install Preview")
	lines = append(lines, fmt.Sprintf("%d installs  ·  %d upgrades  ·  %d reinstalls  ·  %d removals", installs, upgrades, reinstalls, len(summary.Remove)))
	lines = append(lines, "")
	lines = append(lines, "✓ Dependencies resolved")
	lines = append(lines, "✓ Conflicts reviewed")
	lines = append(lines, "")
	if len(packages) != 0 {
		for _, pkg := range packages {
			lines = append(lines, fmt.Sprintf("Warning: %s is already installed.. Reinstalling", pkg))
		}
	}
	finalRows := installSummaryRows(summary)
	finalWidths := installSummaryColumnWidths(finalRows)
	finalLines := installSummaryFormatRows(finalRows, finalWidths)
	lines = append(lines, finalLines...)
	lines = append(lines, "")
	lines = append(lines, downloadString, installedString, previousString)
	return strings.Join(lines, "\n")
}

func mib(number int64) string {
	return fmt.Sprintf("%.2f MiB", float64(number)/(1024*1024))
}

func signed(n int64) string {
	if n < 0 {
		return "−" + mib(-n)
	}
	return mib(n)
}

func installSummaryTotals(instSummary InstallSummary) (int64, int64, int64) {
	var downloadSize int64
	var installedSize int64
	var previousSize int64

	for _, row := range instSummary.Add {
		downloadSize += row.DownloadSize
		installedSize += row.InstalledSize
		previousSize += row.PreviousSize
	}

	for _, row := range instSummary.Remove {
		previousSize += row.PreviousSize
	}
	return downloadSize, installedSize, previousSize
}

func installedSummaryCounts(instSummary InstallSummary) (int, int, int, []string) {
	installs := 0
	reinstalls := 0
	upgrades := 0
	packages := []string{}
	for _, row := range instSummary.Add {
		if len(row.InstalledVersion) == 0 {
			installs++
		} else if row.InstalledVersion == row.NewVersion {
			reinstalls++
			packages = append(packages, fmt.Sprintf("%s: %s", row.Name, row.InstalledVersion))
		} else {
			upgrades++
		}
	}
	return installs, upgrades, reinstalls, packages
}

func installSummaryRows(instSummary InstallSummary) [][]string {
	headerRow := []string{"Package", "Installed Version", "New Version", "Net Change", "Download"}
	buildRows := [][]string{}

	buildRows = append(buildRows, headerRow)

	for _, summary := range instSummary.Add {
		packageRow := []string{}

		packageRow = append(packageRow, fmt.Sprintf("%s/%s", summary.Repo, summary.Name))
		currentVer := summary.InstalledVersion
		if currentVer == "" {
			currentVer = "-"
		}
		packageRow = append(packageRow, currentVer)
		packageRow = append(packageRow, summary.NewVersion)
		packageRow = append(packageRow, signed(summary.InstalledSize-summary.PreviousSize))
		packageRow = append(packageRow, mib(summary.DownloadSize))

		buildRows = append(buildRows, packageRow)
	}
	return buildRows
}

func installSummaryColumnWidths(summaryRows [][]string) []int {
	columnWidth := make([]int, len(summaryRows[0]))

	for _, row := range summaryRows {
		for i, cell := range row {
			rowLength := utf8.RuneCountInString(cell)
			if rowLength > columnWidth[i] {
				columnWidth[i] = rowLength
			}
		}
	}
	return columnWidth
}

func installSummaryFormatRows(summaryRows [][]string, columnWidths []int) []string {
	completeLines := []string{}

	for _, row := range summaryRows {
		rowPadding := []string{}
		for i, cell := range row {
			padding := strings.Repeat(" ", columnWidths[i]-utf8.RuneCountInString(cell))
			if i == 0 || i == 1 || i == 2 {
				columnsOneToThree := cell + padding
				rowPadding = append(rowPadding, columnsOneToThree)
			} else {
				columnsFourFive := padding + cell
				rowPadding = append(rowPadding, columnsFourFive)
			}
		}
		joinedPadding := strings.Join(rowPadding, "   ")
		completeLines = append(completeLines, joinedPadding)
	}

	return completeLines
}
