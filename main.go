package main

import (
	"fmt"
	"os"
)

type Package struct {
	Name    string
	Version string
	Desc    string
}

var Packages = []Package{
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
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please provide a package name to search")
		return
	} else {
		SearchTerm := os.Args[2]
		SearchPackages(SearchTerm)
	}
}
