# Getting Started

tasuke uses a primarily Go stack, with a NodeJS browser application for user registration.
Infrastructure is hosted on GCP, with Firestore used for database. This document will go
through the steps to be ready for local development of any part of tasuke. Individual
component READMEs will document the subset of requirements for that component - if not
needing to develop on any part, feel free to skip unnecessary steps.

## Install tools

The three tools required for development are [Go](https://go.dev/dl/), [NodeJS](https://nodejs.org/en/download/package-manager),
and [gcloud CLI](https://cloud.google.com/sdk/docs/install). Note, for frontend development, only NodeJS is
required and for all other development, NodeJS is not required.

### MacOS

For MacOS users, it is recommended to install tools using Homebrew. First [install Homebrew](https://brew.sh/),
then install the tools.

```bash
brew install go google-cloud-sdk node
```

### Ubuntu (including WSL2)

For Ubuntu users, the default packages in apt will generally be quite old. Follow instructions for each tool
for using a PPA to have more modern packages.

- [gcloud CLI](https://cloud.google.com/sdk/docs/install#deb)
- [Go](https://go.dev/wiki/Ubuntu)
- [NodeJS](https://nodejs.org/en/download/package-manager/all#debian-and-ubuntu-based-linux-distributions)

## Setup

Install `pnpm` using `corepack`, which is part of NodeJS.

```bash
corepack enable
```

Refresh node dependencies

```bash
pnpm i
```

Allow accessing GCP resources (Firestore) by logging into gcloud CLI.

```bash
gcloud auth application-default login
```

You will want to make sure to use a GCP project with a `(default)` Firestore database available.

## IDE

For VSCode users, it is strongly recommended to load the repository's workspace settings using
`File > Open Workspace from File > tasuke.code-workspace`. If prompted, install recommended extensions.
Doing this will ensure formatting runs on save in a way that is consistent with CI.
