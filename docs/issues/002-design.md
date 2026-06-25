# Design: transit CLI and MCP architecture

## Product shape

`transit` has two interfaces backed by one Go core package.

```text
Human / shell    -> cmd/transit CLI -> internal/transit
LLM / Hermes MCP -> cmd/transit mcp -> internal/mcp -> internal/transit
```

## Public API dependency

Base URL: `https://api.transit.ls8h.com`

Initial endpoints:

- `GET /api/v1/locations/suggest`
  - params: `q`, `limit`
  - use: station-name resolution and disambiguation
- `GET /api/v1/plan`
  - params: `from`, `to`, `date`, `time`, `type`, `allowModes`, `avoidModes`, `avoidWalk`, `maxTransfers`, `numItineraries`, `via`
  - use: route planning for departure/arrival/first/last train
- `GET /api/v1/stations/{id}/departures`
  - params: `date`, `time`, `limit`
  - use: station departure board

## Time model

The API returns seconds from service-date midnight in the result timezone. Client formatting must not assume `0 <= seconds < 86400`.

Rules:

- `0` -> `00:00`
- `32700` -> `09:05`
- `91800` -> `01:30(+1d)`
- `-1800` -> `23:30(-1d)`

## CLI design

Default output is human-readable. `--json` returns structured data.

Commands:

```text
transit suggest <query> [--limit N] [--json]
transit plan <from> <to> [--date YYYYMMDD] [--time HH:MM] [--type departure|arrival|first|last] [--num N] [--json]
transit departures <station> [--date YYYYMMDD] [--time HH:MM] [--limit N] [--json]
transit mcp
```

Input resolution:

- `geo:<lat>,<lon>` passes through unchanged.
- Feed-qualified IDs pass through unchanged.
- Plain station names are resolved through `locations/suggest`.
- MVP selects the first station; later work should add disambiguation.

## MCP design

Start with stdio MCP because it works naturally with Hermes and desktop MCP clients.

Tools:

### `transit_suggest_location`

Input:

```json
{ "query": "渋谷", "limit": 5 }
```

Output: station candidates with `id`, `name`, `nameKana`, `feedName`, coordinates.

### `transit_plan_route`

Input:

```json
{
  "from": "渋谷",
  "to": "新宿",
  "date": "20260625",
  "time": "09:00",
  "type": "departure",
  "numItineraries": 3
}
```

Output: compact route options with formatted times, duration, transfer count, and legs. Include raw seconds too for downstream processing.

### `transit_station_departures`

Input:

```json
{ "station": "渋谷", "time": "09:00", "limit": 10 }
```

Output: upcoming departures with formatted departure time, route, mode, headsign.

## Error handling

- API non-2xx: include HTTP status and API error body when available.
- Station not found: clear user-facing error.
- Ambiguous station: MVP picks first; future version should return candidate list and ask caller to choose.
- Invalid flags: fail fast with usage.

## Repository layout

```text
cmd/transit/main.go        CLI entrypoint
internal/transit/client.go API client and DTOs
internal/transit/timefmt.go time formatting
internal/mcp/              future MCP wrapper
docs/architecture.md      architecture notes
```

## Verification

Current scaffold was verified with:

```bash
go test ./...
go build ./cmd/transit
./transit suggest 渋谷 --limit 2
./transit plan 渋谷 新宿 --time 09:00 --num 1
```
