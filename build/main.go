package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/curioswitch/go-build"
	"github.com/curioswitch/go-curiostack/tasks"
	"github.com/goyek/x/boot"
)

func main() {
	// We prefer the shorthand of using non relative `go run build` throughout the repo.
	// It also reduces the number of `build` projects we need within the Go workspace.
	// So while it's an inversion of control to configure projects based on directory
	// structure globally here, it seems to be the simplest solution.
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	switch filepath.Base(wd) {
	case "api":
		tasks.DefineAPI()
	case "server":
		tasks.DefineServer()
	}

	build.DefineTasks(build.LocalPackagePrefix("github.com/curioswitch/tasuke"))
	boot.Main()
}
