# tasuke

tasuke is short for the Japanese 手助け (tedasuke), meaning to give a helping hand.
It aims to support OSS developers in giving each other a helping hand by
connecting code reviewers to PRs in otherwise unrelated OSS repositories.

Everyone knows the [lone maintainer](https://xkcd.com/2347/) issue in OSS -
let's see if we can make it at least a little less lonely.

This project is currently under construction. An initial release will have three
components

- [Frontend web client](./frontend/client) for registering as a code reviewer
- [Frontend API server (BFF)](./frontend/server/) for storing registration information such as allowed review load
- [GitHub webhook](./webhook/server) to listen for code review requests from registered repositories
  and match with a code reviewer based on registration settings

## Related repositories

- [go-curiostack](https://github.com/curioswitch/go-curiostack) - the coding stack used by
  the project, intended to be reusable easily for any new production-grade project
- [tasukeinfra](https://github.com/curioswitch/tasukeinfra) - the Terraform configuration
  for all the infrastructure for tasuke, on top of some bootstrapping in
  [curioinfra](https://github.com/curioswitch/curioinfra).
