package main

import (
	"github.com/curioswitch/go-build"
	"github.com/goyek/x/boot"
)

func main() {
	build.DefineTasks(build.ExcludeTasks("lint-go", "format-go"))
	boot.Main()
}
