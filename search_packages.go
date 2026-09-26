package main

import "fmt"

// SearchPackages -  Real search func but with placeholder code to test functionality before actually searching the repos
func SearchPackages(SearchTerm string) {
	for i, packageName := range Packages {
		if packageName.Name == SearchTerm {
			fmt.Printf("Name: %s\nVersion: %s\nDescription: %s\n", Packages[i].Name, Packages[i].Version, Packages[i].Desc)
		}
	}
}
