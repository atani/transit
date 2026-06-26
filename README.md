# transit

`transit` is a small Go CLI, with an MCP server planned, for querying Japan transit routes and station departures via the public Transit API at <https://api.transit.ls8h.com/>.

## Status

Initial scaffold. The first milestone focuses on a fast human CLI; the MCP server will reuse the same core client and formatter.

## Install

### Homebrew (macOS / Linux)

```bash
brew install atani/tap/transit
```

### winget (Windows)

```powershell
winget install atani.transit
```

### From source

```bash
go install github.com/atani/transit/cmd/transit@latest
```

## Local usage

```bash
go run ./cmd/transit suggest 渋谷
go run ./cmd/transit plan 渋谷 新宿 --type departure --time 09:00
go run ./cmd/transit departures 渋谷 --limit 5
```

The CLI accepts station names and resolves them through `/api/v1/locations/suggest`. Advanced users may pass API endpoints directly, such as `geo:35.681,139.767` or a feed-qualified station ID.

## Planned commands

```text
transit suggest <query>
transit plan <from> <to> [--date YYYYMMDD] [--time HH:MM] [--type departure|arrival|first|last]
transit departures <station> [--date YYYYMMDD] [--time HH:MM] [--limit N]
transit mcp
```

## API notes

Transit API times are seconds from the service-date midnight in the result timezone. Values may exceed 86400 or be negative, so this tool formats them on the client side, including day offsets.

Routing is per-feed: the API plans within a single operator's feed and does not compute transfers between operators. When you pass plain station names, `plan` resolves them to a pair of stations that share a feed so a route can be found. If the origin and destination only exist on different operators (for example a 西鉄 station and a 新幹線 station), no through route exists and the result is a single walking leg; the CLI prints a note in that case.
