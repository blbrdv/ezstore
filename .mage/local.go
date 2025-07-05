//go:build mage

package main

import "github.com/magefile/mage/mg"

// Retest run "lint", "check" and "test" steps.
// Only for local development.
//
//goland:noinspection GoUnusedExportedFunction
func Retest() {
	mg.Deps(Lint)
	mg.Deps(Check)
	mg.Deps(Test)
}

// Rebuild run "clean" and "build" steps.
// Only for local development.
//
//goland:noinspection GoUnusedExportedFunction
func Rebuild() {
	mg.Deps(Clean)
	mg.Deps(Build)
}

// Repack run "clean", "build" and "pack" steps.
// Only for local development.
//
//goland:noinspection GoUnusedExportedFunction
func Repack() {
	mg.Deps(Clean)
	mg.Deps(Build)
	mg.Deps(Pack)
}

// Compile run "build"to target directory.
//
//goland:noinspection GoUnusedExportedFunction
func Compile(targetDir string) error {
	return build(targetDir)
}
