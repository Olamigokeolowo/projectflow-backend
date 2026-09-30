# ProjectFlow API Route Inventory

## Public (no auth required)
| Method | Path | Purpose |
|---|---|---|
| POST | /api/v1/auth/register | Create account |
| POST | /api/v1/auth/login | Get JWT token |

## Authenticated — Workspaces
| Method | Path | Role Required |
|---|---|---|
| POST | /api/v1/workspaces | Member |
| GET | /api/v1/workspaces | Member |
| POST | /api/v1/workspaces/:id/members | Admin |
| DELETE | /api/v1/workspaces/:id/members/:userId | Admin |
| GET | /api/v1/workspaces/:id/members | Member |
| GET | /api/v1/workspaces/:id/activity | Member |

## Authenticated — Decisions
| Method | Path | Auth Rule |
|---|---|---|
| GET | /api/v1/decisions | Membership |
| GET | /api/v1/decisions/:id | Membership |
| POST | /api/v1/decisions | Membership |
| PATCH | /api/v1/decisions/:id | Ownership |
| DELETE | /api/v1/decisions/:id | Ownership |
| GET | /api/v1/decisions/slow | Membership (demo endpoint for cancellation) |

## Authenticated — Tasks
| Method | Path | Auth Rule |
|---|---|---|
| GET/POST | /api/v1/decisions/:id/tasks | Membership (via decision) |
| PATCH/DELETE | /api/v1/tasks/:id | Membership (via decision) |

## Authenticated — Comments
| Method | Path | Auth Rule |
|---|---|---|
| GET/POST | /api/v1/decisions/:id/comments | Membership |
| GET/POST | /api/v1/tasks/:id/comments | Membership (via task → decision) |
| DELETE | /api/v1/comments/:id | Author only |

## Internal / Ops
| Method | Path | Notes |
|---|---|---|
| GET | /metrics | No auth — flagged as an open gap, restrict at network level before production |

## Rate Limits
| Scope | Limit |
|---|---|
| /api/v1/auth/* | 5 requests / minute per IP |
| Everything else | 100 requests / minute per IP |