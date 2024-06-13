package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"strings"

	"connectrpc.com/connect"
	firebase "firebase.google.com/go/v4"
	"github.com/curioswitch/go-curiostack/otel"
	"github.com/curioswitch/go-curiostack/server"
	docshandler "github.com/curioswitch/go-docs-handler"
	protodocs "github.com/curioswitch/go-docs-handler/plugins/proto"
	"github.com/curioswitch/go-usegcp/middleware/firebaseauth"
	"github.com/go-chi/chi/v5/middleware"

	frontendapi "github.com/curioswitch/tasuke/frontend/api/go"
	"github.com/curioswitch/tasuke/frontend/api/go/frontendapiconnect"
	"github.com/curioswitch/tasuke/frontend/server/internal/config"
	"github.com/curioswitch/tasuke/frontend/server/internal/handler/getuser"
	"github.com/curioswitch/tasuke/frontend/server/internal/handler/saveuser"
	"github.com/curioswitch/tasuke/frontend/server/internal/service"
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

	s.Mux().Use(middleware.Maybe(firebaseauth.NewMiddleware(fbAuth), func(r *http.Request) bool {
		return strings.HasPrefix(r.URL.Path, "/"+frontendapiconnect.FrontendServiceName+"/")
	}))

	getUser := getuser.NewHandler(firestore)
	saveUser := saveuser.NewHandler(firestore)
	svc := service.New(getUser, saveUser)

	fapiPath, fapiHandler := frontendapiconnect.NewFrontendServiceHandler(svc,
		connect.WithInterceptors(otel.ConnectInterceptor()))
	s.Mux().Mount(fapiPath, fapiHandler)

	docs, err := docshandler.New(
		protodocs.NewPlugin(
			frontendapiconnect.FrontendServiceName,
			protodocs.WithExampleRequests(
				frontendapiconnect.FrontendServiceGetUserProcedure,
				&frontendapi.GetUserRequest{},
			),
			protodocs.WithExampleRequests(
				frontendapiconnect.FrontendServiceSaveUserProcedure,
				&frontendapi.SaveUserRequest{
					User: &frontendapi.User{
						ProgrammingLanguageIds: []uint32{
							132, // golang
						},
						MaxOpenReviews: 5,
					},
				},
			),
		),
		docshandler.WithInjectedScriptSupplier(func() string {
			script := (`
			function include(url) {
				return new Promise((resolve, reject) => {
					var script = document.createElement('script');
					script.type = 'text/javascript';
					script.src = url;

					script.onload = function() {
						resolve({ script });
					};

					document.getElementsByTagName('head')[0].appendChild(script);
				});
			}

			async function loadScripts() {
				await include("https://alpha.tasuke.dev/__/firebase/8.10.1/firebase-app.js");
				await include("https://alpha.tasuke.dev/__/firebase/8.10.1/firebase-auth.js");
				await include("https://alpha.tasuke.dev/__/firebase/init.js");
				firebase.auth();
			}
			loadScripts();

			async function getAuthorization() {
				const token = await firebase.auth().currentUser.getIdToken();
				return {"Authorization": "Bearer " + token};
			}
			window.armeria.registerHeaderProvider(getAuthorization);
			`)
			return script
		}))
	if err != nil {
		return fmt.Errorf("main: create docs handler: %w", err)
	}
	s.Mux().Handle("/internal/docs/*", http.StripPrefix("/internal/docs", docs))

	return s.Start(ctx)
}
