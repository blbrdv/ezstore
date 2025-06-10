package main

import "github.com/magefile/mage/mg"

// Retest run "check" and "test" steps.
// Only for local development.
func Retest() {
	mg.Deps(Check)
	mg.Deps(Test)
}

// Rebuild run "clean" and "build" steps.
// Only for local development.
func Rebuild() {
	mg.Deps(Clean)
	mg.Deps(Build)
}

// Repack run "clean", "build" and "pack" steps.
// Only for local development.
func Repack() {
	mg.Deps(Clean)
	mg.Deps(Build)
	mg.Deps(Pack)
}
