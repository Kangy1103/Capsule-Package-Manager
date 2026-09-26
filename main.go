package main

import (
	"fmt"
	"os"
)

type Package struct {
	Name    string
	Version string
	Desc    string
	Repo    string
}

/* var Packages = []Package{
	{
		Name:    "neovim",
		Version: "0.11.4-2",
		Desc:    "Fork of Vim aiming to improve extensibility and usability",
	},
	{
		Name:    "vim",
		Version: "9.1.1234-1",
		Desc:    "Vi Improved, a highly configurable text editor",
	},
	{
		Name:    "nano",
		Version: "8.6-1",
		Desc:    "Pico clone with enhancements",
	},
	{
		Name:    "git",
		Version: "2.51.0-1",
		Desc:    "The fast distributed version control system",
	},
} */

// Switch statement needs optimising when more commands are added
// For exmaple, SearchTerm coming from elsewhere instead of being read for every case
// Will also need an error check for os.Args[1]. Possible showing a help message is "capsule" is ran on it's own
func main() {
	PackageInfo, err := FetchPackages()
	if err != nil {
		fmt.Println(err)
		return
	}
	switch os.Args[1] {
	case "search":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a package name to search")
			return
		}
		SearchTerm := os.Args[2]
		SearchPackages(SearchTerm, PackageInfo)
	}
}
