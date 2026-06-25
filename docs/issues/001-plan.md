# Plan: Go transit CLI first, MCP server second

## Summary

Build `transit` as a Go command-line tool for Japan transit routes using `https://api.transit.ls8h.com`, then expose the same core logic through an MCP stdio server.

## Why this order

1. CLI gives immediate value for humans and shell scripts.
2. The API client, station resolver, and formatter can be reused by MCP.
3. MCP should stay thin and return structured tool results rather than duplicating transport/query logic.

## Milestones

### M1: Human CLI MVP

- [x] Create Go module and repository scaffold
- [x] Implement `suggest` command using `/api/v1/locations/suggest`
- [x] Implement `plan` command using `/api/v1/plan`
- [x] Implement `departures` command using `/api/v1/stations/{id}/departures`
- [x] Format service seconds as `HH:MM`, including `(+1d)` / `(-1d)`
- [x] Add `--json` for scripting/agent use
- [x] Add basic unit test for time formatting
- [ ] Add GitHub Actions for `go test ./...` and `go build ./cmd/transit`
- [ ] Improve route formatting, especially fare/platform/transfer display

### M2: Robust CLI

- [ ] Add station disambiguation instead of always selecting the top suggestion
- [ ] Support `geo:<lat>,<lon>` and feed-qualified station IDs explicitly in docs/tests
- [ ] Add `--allow-modes`, `--avoid-modes`, `--avoid-walk`, `--max-transfers`
- [ ] Add shell completion if using a richer CLI framework later
- [ ] Add snapshot/golden tests for output formatting

### M3: MCP server

- [ ] Add `transit mcp` stdio server
- [ ] Expose `transit_suggest_location`
- [ ] Expose `transit_plan_route`
- [ ] Expose `transit_station_departures`
- [ ] Return compact structured JSON suitable for LLM tool results
- [ ] Document Hermes MCP registration command

## Acceptance criteria

```bash
go test ./...
go build ./cmd/transit
./transit suggest 渋谷 --limit 2
./transit plan 渋谷 新宿 --time 09:00 --num 1
```

All commands should complete successfully and produce useful output.
