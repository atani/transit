# Architecture

## Goal

Build a Go-based transit assistant that starts as a fast command-line tool and later exposes the same capability as an MCP server for Hermes/Claude/other agents.

## API

Base URL: `https://api.transit.ls8h.com`

Initial endpoints:

- `GET /api/v1/locations/suggest?q=<query>&limit=<n>`
- `GET /api/v1/plan?from=<id-or-geo>&to=<id-or-geo>&date=<YYYYMMDD>&time=<HH:MM>&type=<departure|arrival|first|last>`
- `GET /api/v1/stations/{id}/departures?date=<YYYYMMDD>&time=<HH:MM>&limit=<n>`

## Components

```text
cmd/transit             CLI entrypoint
internal/transit/client API client and DTOs
internal/transit/format user-facing formatting helpers
internal/mcp            future MCP stdio server wrapper
```

The MCP layer must not duplicate API logic. It should call `internal/transit` and provide JSON-friendly tool results.

## Time handling

Transit API returns times as seconds from service-date midnight in the response timezone. Client formatting must support:

- normal same-day values, e.g. `09:05`
- after-midnight service, e.g. `01:30(+1d)`
- negative previous-day service, e.g. `23:30(-1d)`

## CLI contract

Human-readable output is default. `--json` returns structured API-shaped data for scripting and agents.

## MCP design

Planned tools:

- `transit_suggest_location(query, limit)`
- `transit_plan_route(from, to, date?, time?, type?, numItineraries?)`
- `transit_station_departures(station, date?, time?, limit?)`

The server should use stdio first because Hermes and most MCP clients support it naturally. HTTP can be a later feature if needed.
