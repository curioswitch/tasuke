package main

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/curioswitch/go-build"
	"github.com/curioswitch/go-curiostack/tasks"
	"github.com/goyek/goyek/v2"
	"github.com/goyek/x/boot"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func runNgrok(a *goyek.A) {
	var opts []config.HTTPEndpointOption
	if d := os.Getenv("NGROK_DOMAIN"); d != "" {
		opts = append(opts, config.WithDomain(d))
	}
	lis, err := ngrok.Listen(a.Context(),
		config.HTTPEndpoint(opts...),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		a.Fatalf("Failed to start ngrok: %v", err)
	}

	u, _ := url.Parse("http://localhost:8080")
	p := httputil.NewSingleHostReverseProxy(u)

	a.Logf("Starting ngrok on %s", lis.Addr())
	if err := http.Serve(lis, p); !errors.Is(err, http.ErrServerClosed) { //nolint:gosec // This is a development tool.
		a.Fatalf("Failed to start ngrok: %v", err)
	}
}

func main() {
	tasks.DefineServer()

	build.DefineTasks()

	goyek.Define(goyek.Task{
		Name:  "ngrok",
		Usage: "Starts ngrok to expose the local webhook server to the internet.",
		Action: func(a *goyek.A) {
			runNgrok(a)
		},
	})

	boot.Main()
}
