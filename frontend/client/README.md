# Frontend Client

The web client for registering as a code reviewer on tasuke.

The latest commit is always visible on [https://alpha.tasuke.dev](https://alpha.tasuke.dev).

## Development

This project is a standard web frontend project, so the only hard requirement is [NodeJS](https://nodejs.org/en/download/package-manager).
We use pnpm for package management which must be available. It is generally recommended to use [corepack](https://nodejs.org/api/corepack.html#enabling-the-feature)
for it.

Then, install dependencies and run the dev server by executing the following in this folder. The dev server defaults
to accessing the alpha API server and should be all that is needed for frontend development.

```bash
pnpm install
pnpm run dev
```

VSCode users should open the repository itself as a workspace using `File > Open workspace from file`. If you install
all recommended extensions, formatters and linters will be set up for easy development. You will also find a launch
configuration for running the Frontend Client.
