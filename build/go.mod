module github.com/curioswitch/tasuke/build

go 1.22

require (
	github.com/cli/go-gh/v2 v2.10.0
	github.com/curioswitch/go-build v0.0.0-20240612031439-642962f6d0f9
	github.com/curioswitch/tasuke/common v0.0.0
	github.com/google/go-github/v62 v62.0.0
	github.com/goyek/goyek/v2 v2.2.0
	github.com/goyek/x v0.2.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/cli/safeexec v1.0.0 // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-shellwords v1.0.12 // indirect
	golang.org/x/sys v0.21.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/curioswitch/tasuke/common => ../common/go
