package main

import "fmt"

// SearchPackages -  Real search func but with placeholder code to test functionality before actually searching the repos
func SearchPackages(SearchTerm string, PackageInfo []Package) {
	for i, packageName := range PackageInfo {
		if packageName.Name == SearchTerm {
			fmt.Printf("Name: %s\nVersion: %s\nRepo: %s\nDescription: %s\n",
				PackageInfo[i].Name,
				PackageInfo[i].Version,
				PackageInfo[i].Repo,
				PackageInfo[i].Desc,
			)
		}
	}
}
