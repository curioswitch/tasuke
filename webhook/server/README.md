# Webhook Server

The webhook that listens for review requests from repositories with the tasuke app installed.

## Development

This project is a standard Go project, so only [Go](https://go.dev/dl/) is required for building and
running the code. The database uses GCP Firestore, so you will need to have a GCP project with a
default Firestore database initialized. Make sure you have logged-in to the [gcloud CLI](https://cloud.google.com/sdk/docs/install)
and then you can run the server on http://localhost:8080.

Because this is a webhook, we must also prepare it to be exposed to a GitHub app to listen to events. We
use ngrok for exposing the server. If you don't have one already, register for a free [ngrok account](https://ngrok.com/).
The environment variables `NGROK_AUTHTOKEN` and `NGROK_DOMAIN` should be set to the auth token and domain for your account.
We generally recommend using [direnv](https://github.com/direnv/direnv) and a `.envrc` file for managing local secrets
(this repository `.gitignore`s `.envrc`).

You will also need a GitHub app for development that is tied to your domain. Create a new [GitHub app](https://github.com/settings/apps),
setting the "Webhook URL" to https://your-ngrok-domain.ngrok-free.app/github-webhook (make sure to replace the host with your actual domain)
and set any random string for the webhook secret. Download a private key for the app. Set the environment variables `GITHUB_APPID`,
`GITHUB_SECRET`, and `GITHUB_PRIVATEKEYBASE64` to the application ID of the app, the secret string you defined, and the base64-encoding
of the private key file you downloaded.

Install the app to any test repository you have.

Then, run ngrok in a separate terminal.

```bash
go run ./build ngrok
```

Finally, start the webhook.

```bash
go run ./build start
```

VSCode users should open the repository itself as a workspace using `File > Open workspace from file`. If you install
all recommended extensions, formatters and linters will be set up for easy development. You will also find a launch
configuration for running the Webhook Server + ngrok.
