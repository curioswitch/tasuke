package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"

	"connectrpc.com/connect"
	firebase "firebase.google.com/go/v4"
	"github.com/curioswitch/go-curiostack/otel"
	"github.com/curioswitch/go-curiostack/server"
	docshandler "github.com/curioswitch/go-docs-handler"
	protodocs "github.com/curioswitch/go-docs-handler/plugins/proto"
	"github.com/curioswitch/go-usegcp/middleware/firebaseauth"

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
		log.Fatalf("Failed to create docs handler: %v", err)
	}
	s.Mux().Handle("/internal/docs/*", http.StripPrefix("/internal/docs", docs))

	fbApp, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: conf.Google.Project})
	if err != nil {
		log.Fatalf("Failed to create firebase app: %v", err)
	}

	fbAuth, err := fbApp.Auth(ctx)
	if err != nil {
		log.Fatalf("Failed to create firebase auth client: %v", err)
	}

	firestore, err := fbApp.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to create firestore client: %v", err)
	}
	defer firestore.Close()

	getUser := getuser.NewHandler(firestore)
	saveUser := saveuser.NewHandler(firestore)
	svc := service.New(getUser, saveUser)

	fbMW := firebaseauth.NewMiddleware(fbAuth)
	fapiPath, fapiHandler := frontendapiconnect.NewFrontendServiceHandler(svc,
		connect.WithInterceptors(otel.ConnectInterceptor()))
	fapiHandler = fbMW(fapiHandler)
	s.Mux().Mount(fapiPath, fapiHandler)

	return s.Start(ctx)
}
