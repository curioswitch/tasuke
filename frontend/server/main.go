package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"github.com/curioswitch/go-curiostack/server"
	cstest "github.com/curioswitch/go-curiostack/testutil"
	docshandler "github.com/curioswitch/go-docs-handler"
	protodocs "github.com/curioswitch/go-docs-handler/plugins/proto"
	"github.com/curioswitch/go-usegcp/middleware/firebaseauth"

	frontendapi "github.com/curioswitch/tasuke/frontend/api/go"
	"github.com/curioswitch/tasuke/frontend/api/go/frontendapiconnect"
	"github.com/curioswitch/tasuke/frontend/server/internal/config"
	"github.com/curioswitch/tasuke/frontend/server/internal/handler/saveuser"
	"github.com/curioswitch/tasuke/frontend/server/internal/service"
)

func main() {
	ctx := context.Background()

	conf := config.Load()

	r := server.NewRouter()

	docs, err := docshandler.New(
		protodocs.NewPlugin(
			frontendapiconnect.FrontendServiceName,
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
			token, err := cstest.FirebaseIDToken(context.Background(), "V8yRsCpZJkUfPmxcLI6pKTrx3kf1", "", conf.Google)
			if err != nil {
				log.Printf("Failed to get firebase token: %v", err)
				return ""
			}

			script := fmt.Sprintf(`
			async function getAuthorization() {
				return {"Authorization": "Bearer %s"};
			}
			window.armeria.registerHeaderProvider(getAuthorization);
			`, token)
			return script
		}))
	if err != nil {
		log.Fatalf("Failed to create docs handler: %v", err)
	}
	r.Handle("/internal/docs/*", http.StripPrefix("/internal/docs", docs))

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

	saveUser := saveuser.NewHandler(firestore)
	svc := service.New(saveUser)

	fbMW := firebaseauth.NewMiddleware(fbAuth)
	fapiPath, fapiHandler := frontendapiconnect.NewFrontendServiceHandler(svc)
	fapiHandler = fbMW(fapiHandler)
	r.Mount(fapiPath, fapiHandler)

	srv := server.NewServer(r, conf.Server)

	log.Printf("Starting server on address %v\n", srv.Addr)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Failed to start server: %v", err)
	}
}
