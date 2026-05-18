# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run the server
go run ./cmd/server

# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/messages/...

# Run a single test
go test ./internal/messages/... -run TestService_ListByChatAfter
```

## Architecture

The server is a single HTTP process with all state held in memory (no database). Data is lost on restart.

**Layer structure per domain:**
- `internal/<domain>/service.go` — business logic, owns the domain structs
- `internal/<domain>/errors.go` — sentinel errors for the domain
- `internal/<domain>/id.go` — domain-specific ID type (thin wrapper over `shared/id`)
- `internal/httpserver/<domain>.go` — HTTP handler that calls the service directly

**Request flow:** `httpserver.New` (wires routes) → `Handler` (holds all three services) → domain `Service` → in-memory storage.

**ID pattern:** Each domain defines its own `type ID id.ID` so user/chat/message IDs are distinct types and can't be accidentally swapped. `shared/id` contains the UUID generation and parsing logic shared by all domains.

**Messages pagination:** `Service.ListByChatAfter(chatID, after, limit)` — if `after` is empty, starts from the beginning; if the `after` ID isn't found in the list, it also starts from the beginning. Default limit is 50.

**Testing pattern:** Services are tested with an in-package `fakeRepository` that implements the `Repository` interface. Only `messages` has a `Repository` interface + `MemoryRepository` implementation; `users` and `chats` embed their maps directly in the service.
