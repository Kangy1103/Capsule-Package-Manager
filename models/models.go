// Package models is basically where structs live that are used across other packages
package models

import "github.com/Jguer/go-alpm/v2"

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

type PackageInfo struct {
	DetailedPackageInfo
	Source alpm.IPackage
}
