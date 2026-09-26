package main

import (
	"fmt"
	"strings"

	"github.com/Jguer/go-alpm/v2"
	pacmanconf "github.com/Morganamilo/go-pacmanconf"
)

func FetchPackages() ([]DetailedPackageInfo, error) {
	PackageInfo := []DetailedPackageInfo{}

	// This part initialises a libalpm handle for the pacman database path
	handle, err := alpm.Initialize("/", "/var/lib/pacman/")
	if err != nil {
		fmt.Printf("This didn't work: %s\n", err)
	}

	// This is parsing the pacman.conf using go-pacmanconf
	conf, _, err := pacmanconf.ParseFile("/etc/pacman.conf")
	if err != nil {
		return []DetailedPackageInfo{}, fmt.Errorf("pacman.conf not found: %s", err)
	}

	// This is looping over the repos to register the defined repos in pacman.conf for libalpm
	for _, repo := range conf.Repos {
		db, err := handle.RegisterSyncDB(repo.Name, alpm.SigUseDefault)
		if err != nil {
			fmt.Printf("Could not register %s: %s\n", repo.Name, err)
			continue
		}

		// Gives libalpm the repo list
		db.SetServers(repo.Servers)
	}

	dbs, err := handle.SyncDBs()
	if err != nil {
		fmt.Printf("handle.SyncDBs failed: %s\n", err)
	}
	databases := dbs.Slice()

	// This loops all repos for the package cache
	for _, db := range databases {
		packages := db.PkgCache()

		// Formats the packages for Capsule to understand: IPackage to BasicPackageInfo
		for _, pkg := range packages.Slice() {
			PackageInfo = append(PackageInfo, DetailedPackageInfo{
				BasicPackageInfo: BasicPackageInfo{
					Name:    pkg.Name(),
					Version: pkg.Version(),
					Desc:    pkg.Description(),
					Repo:    db.Name(),
				},
				Architecture:  pkg.Architecture(),
				URL:           pkg.URL(),
				Licenses:      strings.Join(pkg.Licenses().Slice(), "  "),
				Groups:        strings.Join(pkg.Groups().Slice(), "  "),
				Provides:      formatDependencies(pkg.Provides()),
				DependsOn:     formatDependencies(pkg.Depends()),
				OptionalDeps:  formatDependencies(pkg.OptionalDepends()),
				ConflictsWith: formatDependencies(pkg.Conflicts()),
				Replaces:      formatDependencies(pkg.Replaces()),
				DownloadSize:  formatFileSize(pkg.Size()),
				InstalledSize: formatFileSize(pkg.ISize()),
				Packager:      pkg.Packager(),
				BuildDate:     pkg.BuildDate().Format("Mon 02 Jan 2006 15:04:05 MST"),
				ValidatedBy:   formatValidatedBy(pkg.Validation()),
			})
		}
	}
	return PackageInfo, nil
}

func formatDependencies(deps alpm.IDependList) string {
	depsSlice := deps.Slice()
	depsString := []string{}

	for _, dep := range depsSlice {
		depsString = append(depsString, dep.Name)
	}
	return strings.Join(depsString, ", ")
}

func formatFileSize(fileSize int64) string {
	fileSizeBytes := float64(fileSize)
	fileSizeMiB := fileSizeBytes / 1024 / 1024
	return fmt.Sprintf("%.2f MiB", fileSizeMiB)
}

func formatValidatedBy(validation alpm.Validation) string {
	validatedByString := []string{}
	if validation&alpm.ValidationSHA256Sum != 0 {
		validatedByString = append(validatedByString, "SHA-256 Sum")
	}
	if validation&alpm.ValidationMD5Sum != 0 {
		validatedByString = append(validatedByString, "MD5 Sum")
	}
	if validation&alpm.ValidationSignature != 0 {
		validatedByString = append(validatedByString, "Signature")
	}
	if validation&alpm.ValidationNone != 0 {
		validatedByString = append(validatedByString, "None")
	}
	return strings.Join(validatedByString, ", ")
}
