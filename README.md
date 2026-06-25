# transit

`transit` is a small Go CLI, with an MCP server planned, for querying Japan transit routes and station departures via the public Transit API at <https://api.transit.ls8h.com/>.

## Status

Initial scaffold. The first milestone focuses on a fast human CLI; the MCP server will reuse the same core client and formatter.

## Install from source

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
