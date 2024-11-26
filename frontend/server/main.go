package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"github.com/curioswitch/go-curiostack/server"
	"github.com/curioswitch/go-usegcp/middleware/firebaseauth"
	"github.com/go-chi/chi/v5/middleware"

	frontendapi "github.com/curioswitch/tasuke/frontend/api/go"
	"github.com/curioswitch/tasuke/frontend/api/go/frontendapiconnect"
	"github.com/curioswitch/tasuke/frontend/server/internal/config"
	"github.com/curioswitch/tasuke/frontend/server/internal/handler/getuser"
	"github.com/curioswitch/tasuke/frontend/server/internal/handler/saveuser"
)

// e2e-test1@curioswitch.org
const e2eTest1UID = "V8yRsCpZJkUfPmxcLI6pKTrx3kf1"

var confFiles embed.FS // Currently empty

func main() {
	os.Exit(server.Main(&config.Config{}, confFiles, setupServer))
}

func setupServer(ctx context.Context, conf *config.Config, s *server.Server) error {
	fbApp, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: conf.Google.Project})
	if err != nil {
		return fmt.Errorf("main: create firebase app: %w", err)
	}

	fbAuth, err := fbApp.Auth(ctx)
	if err != nil {
		return fmt.Errorf("main: create firebase auth client: %w", err)
	}

	firestore, err := fbApp.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("main: create firestore client: %w", err)
	}
	defer firestore.Close()

	server.Mux(s).Use(middleware.Maybe(firebaseauth.NewMiddleware(fbAuth), func(r *http.Request) bool {
		return strings.HasPrefix(r.URL.Path, "/"+frontendapiconnect.FrontendServiceName+"/")
	}))

	getUser := getuser.NewHandler(firestore)
	saveUser := saveuser.NewHandler(firestore)

	server.HandleConnectUnary(s,
		frontendapiconnect.FrontendServiceGetUserProcedure,
		getUser.GetUser,
		[]*frontendapi.GetUserRequest{
			{},
		},
	)
	server.HandleConnectUnary(s,
		frontendapiconnect.FrontendServiceSaveUserProcedure,
		saveUser.SaveUser,
		[]*frontendapi.SaveUserRequest{
			{
				User: &frontendapi.User{
					ProgrammingLanguageIds: []uint32{
						132, // golang
					},
					MaxOpenReviews: 5,
				},
			},
		},
	)

	server.EnableDocsFirebaseAuth(s, "alpha.tasuke.dev")

	return server.Start(ctx, s)
}
