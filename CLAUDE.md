# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GoFlex is a Go library and CLI tool for interacting with the Plex Media Server API. The primary use case is intelligent playlist randomization that removes watched episodes and refills playlists automatically. It started as a playlist randomizer but has expanded to include general Plex API operations.

## Environment Setup

Required environment variables:
```bash
PLEX_URL=http://<IP>:32400
PLEX_TOKEN=<TOKEN>
```

## Common Commands

### Building
```bash
# Build the CLI binary
go build -o goflex ./cmd/goflex

# Build using goreleaser (for releases)
goreleaser build --snapshot --clean
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a specific test
go test -v -run TestName ./path/to/package

# Run tests for a specific file
go test -v ./playlist_test.go
```

### Linting
```bash
# Run golangci-lint (uses .golangci.yml config)
golangci-lint run

# Auto-fix issues where possible
golangci-lint run --fix
```

### Running the CLI
```bash
# Get playlists
./goflex get playlists

# Randomize playlists using config file(s)
./goflex random config.yaml

# The random command can process multiple configs in parallel
./goflex random config-1.yaml config-2.yaml
```

## Architecture

### Core Design Pattern: Service-Oriented Architecture

The codebase uses a service-oriented pattern where the main `Flex` struct (defined in `plex.go`) acts as a client that aggregates multiple service interfaces:

```go
type Flex struct {
    // Core configuration
    baseURL, token string
    client *http.Client
    cache cache

    // Services
    Playlists      PlaylistService
    Sessions       SessionService
    Media          MediaService
    Server         ServerService
    Shows          ShowService
    Library        LibraryService
    Authentication AuthenticationService
}
```

Each service is defined as an interface with a corresponding `*ServiceOp` implementation that holds a reference back to the parent `Flex` instance.

### Key Components

**plex.go**: Main client initialization using functional options pattern. The `New()` function accepts variadic options (e.g., `WithBaseURL()`, `WithToken()`, `WithFlexConfig()`).

**api.go**: Core HTTP request/response handling. Contains XML and JSON response structs for Plex API endpoints. The `sendRequestType()` method handles caching, content-type negotiation, and HTTP execution.

**Cache System**: In-memory cache (`cache.go`) with TTL support and garbage collection. Cache keys are generated from HTTP request URLs with prefixes. Services can invalidate cache entries by prefix (e.g., when a playlist is modified).

**playlist.go**: The most complex service. Implements the `Randomize()` method which:
1. Gets or creates the target playlist
2. Fetches episodes currently in the playlist
3. Queries viewing history to identify watched episodes
4. Removes watched episodes from the playlist
5. If below refill threshold, adds random unviewed episodes from configured series
6. Calculates sleep duration based on current playback status

**show.go**, **episode.go**, **sessions.go**: Domain-specific services for TV shows, episodes, and viewing sessions. These interact with different Plex API endpoints and handle the XML/JSON response parsing.

**library.go**, **server.go**, **media.go**: Supporting services for library management, server info, and media operations.

### CLI Structure

Located in `cmd/goflex/`:
- `cmd/root.go`: Root command with global flags (`--verbose`, `--gc-interval`)
- `cmd/random.go`: Main randomization command that loads YAML configs and runs playlist randomizers concurrently using `errgroup`
- `cmd/get_*.go`: Various "get" subcommands for retrieving Plex resources
- `cmd/create_*.go`, `cmd/delete_*.go`: CRUD operations for playlists and other resources

The CLI uses Cobra for command structure and Charmbracelet's log library for pretty terminal output.

### Important Type Patterns

**Custom string types**: `PlaylistTitle`, `ShowTitle`, `SeasonNumber`, `EpisodeNumber` are typed aliases for better type safety and self-documenting code.

**Functional options**: Both the library (`Flex`) and key functions (`NewRandomizeRequest`) use the functional options pattern for configuration.

**EpisodeList operations**: The `EpisodeList` type (slice of `Episode`) has a `Subtract()` method that returns episodes remaining after removing a subset, used for filtering watched episodes.

### Testing

Tests use the `testdata/` directory with recorded XML/JSON responses from Plex API. This allows testing without a live Plex server. Test files are co-located with their implementation files (e.g., `playlist.go` and `playlist_test.go`).

### Response Caching Strategy

API responses are cached with different TTLs based on volatility:
- Playlists list: 10 minutes
- Playlist episodes: 60 minutes
- Seasons: 1 hour

Cache is invalidated explicitly when mutations occur (create/delete playlist, remove episodes, etc.) using prefix-based deletion.

### Randomization Loop

The `random` command runs an infinite loop per playlist that:
1. Executes randomization logic
2. Calculates next check time based on current viewing progress
3. Sleeps until next check
4. Repeats

Multiple randomizers run concurrently using `errgroup`, allowing a single process to manage multiple playlists across potentially multiple Plex servers.

## Release Process

Releases are automated via GitHub Actions (`.github/workflows/release.yml`) using GoReleaser when tags matching `v*` are pushed. The GoReleaser config (`.goreleaser.yaml`) builds for Linux, Windows, and Darwin.
