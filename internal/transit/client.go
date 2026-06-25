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
		return "", "", fmt.Errorf("station not found: %s", input)
	}
	return res.Stations[0].ID, res.Stations[0].Name, nil
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
