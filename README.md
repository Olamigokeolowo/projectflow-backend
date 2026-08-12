# ProjectFlow

A backend service exploring real-world Go/Gin engineering patterns:
routing, middleware, layered architecture, auth, caching, and observability.

ProjectFlow is a backend service for tracking decisions and the tasks that come out of them. Think of it as the engine behind something like "we decided to do X, here's why and here's what needs to happen because of it". A structured way to record choices a team makes and follow through on them.

Strip away the tech and it's genuinely simple: a decision has a title, a status (draft, proposed, accepted, rejected), an owner, and it can spawn tasks. That's the whole domain. The complexity you're building isn't in what it does — it's in doing it the way production systems actually work: authenticated, validated, layered, observable, resilient to restarts and failures. That's the real point of this project. You're not learning to build a to-do app; you're learning to build a to-do app correctly, the way it'd need to be built if real users and real data depended on it.

## Status
🚧 In progress - building in public, one concept at a time.

## Stack
- Go
- Gin