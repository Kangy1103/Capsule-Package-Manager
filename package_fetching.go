package main

import (
	"fmt"

	"github.com/Jguer/go-alpm/v2"
	pacmanconf "github.com/Morganamilo/go-pacmanconf"
)

func FetchPackages() ([]Package, error) {
	PackageInfo := []Package{}

	// This part initialises a libalpm handle for the pacman database path
	handle, err := alpm.Initialize("/", "/var/lib/pacman/")
	if err != nil {
		fmt.Printf("This didn't work: %s\n", err)
	}

	// This is parsing the pacman.conf using go-pacmanconf
	conf, _, err := pacmanconf.ParseFile("/etc/pacman.conf")
	if err != nil {
		return []Package{}, fmt.Errorf("pacman.conf not found: %s", err)
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

		// Formats the packages for Capsule to understand: IPackage to Package
		for _, pkg := range packages.Slice() {
			PackageInfo = append(PackageInfo, Package{
				Name:    pkg.Name(),
				Version: pkg.Version(),
				Desc:    pkg.Description(),
				Repo:    db.Name(),
			})
		}
	}
	return PackageInfo, nil
}
