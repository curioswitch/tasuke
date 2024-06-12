package main

import (
	"github.com/curioswitch/go-build"
	"github.com/curioswitch/go-curiostack/tasks"
	"github.com/goyek/x/boot"
)

func main() {
	tasks.DefineAPI()
	build.DefineTasks()
	boot.Main()
}
