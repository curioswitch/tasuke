# Frontend Server

The API server for the registration web client.

## Development

This project is a standard Go project, so only [Go](https://go.dev/dl/) is required for building and
running the code. The database uses GCP Firestore, so you will need to have a GCP project with a
default Firestore database initialized. Make sure you have logged-in to the [gcloud CLI](https://cloud.google.com/sdk/docs/install)
and then you can run the server on http://localhost:8080.

```bash
go run ./build start
```

The easiest way to check it is working is to use the debug page at http://localhost:8080/internal/docs/.

VSCode users should open the repository itself as a workspace using `File > Open workspace from file`. If you install
all recommended extensions, formatters and linters will be set up for easy development. You will also find a launch
configuration for running the Frontend Server.
