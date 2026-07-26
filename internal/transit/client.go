package transit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const DefaultBaseURL = "https://api.transit.ls8h.com"

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return &Client{BaseURL: strings.TrimRight(baseURL, "/"), HTTP: &http.Client{Timeout: 15 * time.Second}}
}

type Station struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	NameKana string  `json:"nameKana,omitempty"`
	FeedID   string  `json:"feedId"`
	FeedName string  `json:"feedName"`
	Kind     string  `json:"kind,omitempty"`
	Lat      float64 `json:"lat,omitempty"`
	Lon      float64 `json:"lon,omitempty"`
}

type SuggestResponse struct {
	Stations []Station `json:"stations"`
}

type StopRef struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	PlatformCode string `json:"platformCode,omitempty"`
}

type Leg struct {
	Kind          string  `json:"kind"`
	RouteName     string  `json:"routeName,omitempty"`
	Mode          string  `json:"mode,omitempty"`
	Headsign      string  `json:"headsign,omitempty"`
	From          StopRef `json:"from"`
	To            StopRef `json:"to"`
	DepartureSecs int     `json:"departureSecs"`
	ArrivalSecs   int     `json:"arrivalSecs"`
}

type Journey struct {
	DepartureSecs int   `json:"departureSecs"`
	ArrivalSecs   int   `json:"arrivalSecs"`
	DurationSecs  int   `json:"durationSecs"`
	TransferCount int   `json:"transferCount"`
	Legs          []Leg `json:"legs"`
}

type PlanResponse struct {
	Date     string    `json:"date"`
	Type     string    `json:"type"`
	Timezone string    `json:"timezone"`
	From     StopRef   `json:"from"`
	To       StopRef   `json:"to"`
	Journeys []Journey `json:"journeys"`
}

type Departure struct {
	RouteName     string `json:"routeName"`
	Mode          string `json:"mode"`
	Headsign      string `json:"headsign,omitempty"`
	DepartureSecs int    `json:"departureSecs"`
}

type DeparturesResponse struct {
	StationID  string      `json:"stationId"`
	Date       string      `json:"date"`
	Timezone   string      `json:"timezone"`
	Departures []Departure `json:"departures"`
}

type PlanOptions struct {
	Date, Time, Type string
	NumItineraries   int
}

type DepartureOptions struct {
	Date, Time string
	Limit      int
}

func (c *Client) Suggest(ctx context.Context, query string, limit int) (SuggestResponse, error) {
	var out SuggestResponse
	q := url.Values{}
	q.Set("q", query)
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	return out, c.get(ctx, "/api/v1/locations/suggest", q, &out)
}

func (c *Client) Plan(ctx context.Context, from, to string, opts PlanOptions) (PlanResponse, error) {
	var out PlanResponse
	q := url.Values{}
	q.Set("from", from)
	q.Set("to", to)
	if opts.Date != "" {
		q.Set("date", opts.Date)
	}
	if opts.Time != "" {
		q.Set("time", opts.Time)
	}
	if opts.Type != "" {
		q.Set("type", opts.Type)
	}
	if opts.NumItineraries > 0 {
		q.Set("numItineraries", strconv.Itoa(opts.NumItineraries))
	}
	return out, c.get(ctx, "/api/v1/plan", q, &out)
}

func (c *Client) Departures(ctx context.Context, stationID string, opts DepartureOptions) (DeparturesResponse, error) {
	var out DeparturesResponse
	q := url.Values{}
	if opts.Date != "" {
		q.Set("date", opts.Date)
	}
	if opts.Time != "" {
		q.Set("time", opts.Time)
	}
	if opts.Limit > 0 {
		q.Set("limit", strconv.Itoa(opts.Limit))
	}
	path := "/api/v1/stations/" + url.PathEscape(stationID) + "/departures"
	return out, c.get(ctx, path, q, &out)
}

func (c *Client) ResolveStation(ctx context.Context, input string) (string, string, error) {
	if strings.HasPrefix(input, "geo:") || strings.Contains(input, ":") {
		return input, input, nil
	}
	res, err := c.Suggest(ctx, input, 1)
	if err != nil {
		return "", "", err
	}
	if len(res.Stations) == 0 {
		return "", "", fmt.Errorf("駅が見つかりません: %s", input)
	}
	return res.Stations[0].ID, res.Stations[0].Name, nil
}

// candidates returns station suggestions for a query, or a single synthetic
// station for passthrough inputs (geo: coordinates or feed-qualified IDs).
func (c *Client) candidates(ctx context.Context, input string) ([]Station, error) {
	if strings.Contains(input, ":") {
		return []Station{{ID: input, Name: input}}, nil
	}
	res, err := c.Suggest(ctx, input, 10)
	if err != nil {
		return nil, err
	}
	if len(res.Stations) == 0 {
		return nil, fmt.Errorf("駅が見つかりません: %s", input)
	}
	return res.Stations, nil
}

// ResolveStationPair resolves from and to into the pair of stops most likely to
// yield a sensible itinerary. A plain name such as 押上 or 渋谷 matches both rail
// stations and bus stops across many feeds, and the choice drives the result.
func (c *Client) ResolveStationPair(ctx context.Context, from, to string) (Station, Station, error) {
	fromCands, err := c.candidates(ctx, from)
	if err != nil {
		return Station{}, Station{}, err
	}
	toCands, err := c.candidates(ctx, to)
	if err != nil {
		return Station{}, Station{}, err
	}
	f, t := pickFeedPair(fromCands, toCands)
	return f, t, nil
}

// stationKind marks a suggestion as a rail station rather than a bus stop.
const stationKind = "station"

// railStations returns only the rail-station candidates, preserving rank.
func railStations(cands []Station) []Station {
	out := make([]Station, 0, len(cands))
	for _, s := range cands {
		if s.Kind == stationKind {
			out = append(out, s)
		}
	}
	return out
}

// sharedFeed returns the highest-ranked from/to pair that shares a feed.
func sharedFeed(from, to []Station) (Station, Station, bool) {
	for _, f := range from {
		for _, t := range to {
			if f.FeedID != "" && f.FeedID == t.FeedID {
				return f, t, true
			}
		}
	}
	return Station{}, Station{}, false
}

// pickFeedPair chooses the from/to pair to plan with, in descending preference:
//
//  1. two rail stations sharing a feed — a direct ride with no transfer between
//     operators, e.g. 渋谷→新宿 on the 埼京線;
//  2. two rail stations on different feeds — the API routes across feeds, so
//     this is fine, e.g. 和白 (西鉄) → 博多 (JR);
//  3. anything sharing a feed, then the top candidates, for places that only
//     exist as bus stops.
//
// Preferring rail stations matters because a name like 押上 also matches bus
// stops: resolving to those yields an all-bus itinerary that can span a day.
func pickFeedPair(from, to []Station) (Station, Station) {
	fromRail, toRail := railStations(from), railStations(to)
	if f, t, ok := sharedFeed(fromRail, toRail); ok {
		return f, t
	}
	if len(fromRail) > 0 && len(toRail) > 0 {
		return fromRail[0], toRail[0]
	}
	if f, t, ok := sharedFeed(from, to); ok {
		return f, t
	}
	return from[0], to[0]
}

func (c *Client) get(ctx context.Context, path string, q url.Values, out any) error {
	u := c.BaseURL + path
	if len(q) > 0 {
		u += "?" + q.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "transit/0.1")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// The error body may be JSON or a plain-text/HTML page; read it raw
		// with a size cap so a huge response can't exhaust memory.
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
		if msg := strings.TrimSpace(string(raw)); msg != "" {
			return fmt.Errorf("transit api %s: %s", resp.Status, msg)
		}
		return fmt.Errorf("transit api %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
