# Notable rationale of tasuke

Below documents some design decisions in tasuke. They aim to provide some background and potential food
for thought if a reader is designing a system and in the same situation. These rarely push for a
"correct answer" (if something was obviously correct, it likely wouldn't need to be documented here)
but are approaches taken with pros/cons that fit with the needs or philosophies of the main developers.

## Use firebase hosting for API ingress

With GCP, it is generally preferred to use a [load balancer](https://cloud.google.com/load-balancing)
for ingress to a service hosted with Cloud Run because it enables useful features such as
CDN and IAP protection of internal endpoints, however a load balancer has a fixed running
cost of ~$20 per month.

We use firebase hosting's ability to proxy to Cloud Run instead only because it has no cost and
this project aims to be as lean as possible to run without funding.

## Use Cloud Firestore for database

Firestore can be tricky to use due to its limited querying model, compared to a full database such
as PostgreSQL. We use Firestore though because it has no fixed running cost and a generous free
tier, as this project aims to be as lean as possible to run without funding.

## Implement frontend API server in Go

It is common for frontends to implement their server (i.e. BFF) in NodeJS to have consistency with
the client code, with NextJS being particularly popular. We prefer the development efficiency of Go
and want to use it as much as possible - notably, we can share utilities for server-side development
between the API server and webhook. By using an RPC framework (connect), it is trivial to call into
the Go server side from the browser even though the languages don't match, so we feel this is more
maintainable than doing everything in NodeJS.

## Implement frontend client with Vike

Vite is a popular choice for bundling browser applications without a server-side component. We also
want to be able to prerender pages before deployment, which Vike makes very simple.

## Require review requesters to be registered reviewers

Initially, we don't know how many reviewers will be available or what abuse models are possible
through the service. We require requesters to have an account with at least one review allowed
while understanding how the service is used to increase the chance there are available reviewers.

In the long term, it could be nice if reviewers can help learners that are getting started but
can't do reviews, and to enable that, in the future this restriction would need to be lifted.
