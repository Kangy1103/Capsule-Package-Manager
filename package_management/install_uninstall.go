// Package packagemanagement is where the install/uninstall code lives
package packagemanagement

import (
	"fmt"
	"strings"

	"capsule-package-manager/models"

	"charm.land/log/v2"
	"github.com/Jguer/go-alpm/v2"
)

func InstallPackage(searchTerm string, handle *alpm.Handle, pkgs []models.PackageInfo) {
	var installPkg models.PackageInfo
	var found bool

	for _, pkg := range pkgs {
		if searchTerm != pkg.Name {
			continue
		}
		installPkg = pkg
		found = true
		break
	}
	if !found {
		log.Error("package not found", "pkgmgmnt-err", searchTerm)
		fmt.Printf("package not found: %v", searchTerm)
		return
	}

	if err := handle.TransInit(0); err != nil {
		log.Error("unable to initialise package transaction", "pkgmgmnt-err", err)
		return
	}
	defer handle.TransRelease()

	// Adding pkgs to be installed
	log.Info("Adding packages to install")
	if err := handle.AddPkg(installPkg.Source); err != nil {
		log.Error("unable to add pkg", "pkgmgmnt-err", err)
		return
	}

	conflictCtx := context{
		transactionHandle: handle,
	}
	// Question Handler
	log.Info("Calling questionHandler()")
	handle.SetQuestionCallback(questionHandler, &conflictCtx)

	if err := handle.TransPrepare(); err != nil {
		log.Error("dependency check failed", "pkgmgmnt-err", err)
		return
	}

	// List in console so user can see what's going to change
	fmt.Println("Packages to install:")
	for _, pkg := range handle.TransGetAdd().Slice() {
		fmt.Printf(" %s %s\n", pkg.Name(), pkg.Version())
	}
	fmt.Println(" ")
	removalList := handle.TransGetRemove().Slice()
	if len(removalList) != 0 {
		fmt.Println("Package to remove:")
		for _, pkg := range removalList {
			fmt.Printf(" %s %s\n", pkg.Name(), pkg.Version())
		}
		fmt.Println(" ")
	}
	fmt.Print("Do you want to proceed? [y/N] ")
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		log.Error("could not read response", "pkgmgmnt-err", err)
		fmt.Println("Could not read response")
		return
	}
	lowerResponse := strings.ToLower(response)
	if lowerResponse != "y" && lowerResponse != "yes" {
		fmt.Println("Installation cancelled")
		return
	}

	// Commit the install
	log.Info("Commiting packages to install")
	if err := handle.TransCommit(); err != nil {
		log.Error("installation failed", "pkgmgmnt-err", err)
		fmt.Printf("Installation failed: %v", err)
		return
	}
	log.Info("Installation complete!")
	fmt.Println("Installation complete!")
}

func UninstallPackage() {}

type context struct {
	transactionHandle *alpm.Handle
	checkComplete     bool
}

func questionHandler(ctx any, question alpm.QuestionAny) {
	conflictCtx := ctx.(*context)
	handle := conflictCtx.transactionHandle
	if question.Type() == alpm.QuestionTypeConflictPkg {
		question.SetAnswer(false)
		if conflictCtx.checkComplete {
			question.SetAnswer(true)
			return
		}
		pkgs := handle.TransGetAdd()
		localDB, err := handle.LocalDB()
		if err != nil {
			log.Error("unable to fetch local DB", "pkgmgmnt-err", err)
			return
		}
		pkgCache := localDB.PkgCache()
		approvedRemovals := make(map[string]bool)
		excludePkgs := make(map[string]bool)

		// Incoming packages
		for _, pkg := range pkgs.Slice() {
			excludePkgs[pkg.Name()] = true
		}

		// Check incoming packages
		for _, pkg := range pkgs.Slice() {
			pkgConflicts := pkg.Conflicts()
			for _, pkgRule := range pkgConflicts.Slice() {
				for _, conflictingPkg := range pkgCache.Slice() {
					if ok := excludePkgs[conflictingPkg.Name()]; ok {
						continue
					}
					if !matchesConflictRule(conflictingPkg, pkgRule) {
						continue
					}
					if _, ok := approvedRemovals[conflictingPkg.Name()]; ok {
						continue
					}
					fmt.Printf(".: %s-%s and %s-%s conflict. Remove %s? [y/N] :.", pkg.Name(), pkg.Version(), conflictingPkg.Name(), conflictingPkg.Version(), conflictingPkg.Name())

					var response string
					if _, err := fmt.Scanln(&response); err != nil {
						log.Error("could not read response", "pkgmgmnt-err", err)
						fmt.Println("Could not read response")
						return
					}

					if strings.ToLower(response) == "n" || strings.ToLower(response) == "no" {
						log.Error("install cancelled", "pkgmgmnt-err", response)
						question.SetAnswer(false)
						return
					} else if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
						approvedRemovals[conflictingPkg.Name()] = true
					} else {
						return
					}
				}
			}
		}

		// Check installed packages
		for _, pkg := range pkgCache.Slice() {
			if _, ok := excludePkgs[pkg.Name()]; ok {
				continue
			}
			pkgConflicts := pkg.Conflicts()
			for _, pkgRule := range pkgConflicts.Slice() {
				for _, conflictingPkg := range pkgs.Slice() {
					if !matchesConflictRule(conflictingPkg, pkgRule) {
						continue
					}
					if _, ok := approvedRemovals[pkg.Name()]; ok {
						continue
					}
					fmt.Printf(".: %s-%s and %s-%s conflict. Remove %s? [y/N] :.", pkg.Name(), pkg.Version(), conflictingPkg.Name(), conflictingPkg.Version(), pkg.Name())

					var response string
					if _, err := fmt.Scanln(&response); err != nil {
						log.Error("could not read response", "pkgmgmnt-err", err)
						fmt.Println("Could not read response")
						return
					}

					if strings.ToLower(response) == "n" || strings.ToLower(response) == "no" {
						log.Error("install cancelled", "pkgmgmnt-err", response)
						question.SetAnswer(false)
						return
					} else if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
						approvedRemovals[pkg.Name()] = true
					} else {
						return
					}
				}
			}
		}
		conflictCtx.checkComplete = true
		question.SetAnswer(true)
	}
}

func matchesConflictRule(pkg alpm.IPackage, conflictRule alpm.Depend) bool {
	if pkg.Name() == conflictRule.Name {
		if conflictRule.Version == "" {
			return true
		}
		pkgVersion := alpm.VerCmp(pkg.Version(), conflictRule.Version)

		switch conflictRule.Mod {
		case alpm.DepModEq:
			if pkgVersion == 0 {
				return true
			}
		case alpm.DepModGE:
			if pkgVersion >= 0 {
				return true
			}
		case alpm.DepModLE:
			if pkgVersion <= 0 {
				return true
			}
		case alpm.DepModGT:
			if pkgVersion > 0 {
				return true
			}
		case alpm.DepModLT:
			if pkgVersion < 0 {
				return true
			}
		}
	}

	for _, pkgProvision := range pkg.Provides().Slice() {
		if pkgProvision.Name == conflictRule.Name && conflictRule.Version == "" {
			return true
		}
		if pkgProvision.Name != conflictRule.Name || pkgProvision.Version == "" {
			continue
		}
		pkgVer := alpm.VerCmp(pkgProvision.Version, conflictRule.Version)
		switch conflictRule.Mod {
		case alpm.DepModEq:
			if pkgVer == 0 {
				return true
			}
		case alpm.DepModGE:
			if pkgVer >= 0 {
				return true
			}
		case alpm.DepModLE:
			if pkgVer <= 0 {
				return true
			}
		case alpm.DepModGT:
			if pkgVer > 0 {
				return true
			}
		case alpm.DepModLT:
			if pkgVer < 0 {
				return true
			}
		}
	}
	return false
}
