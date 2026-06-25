package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/atani/transit/internal/transit"
)

// version is overridden at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		usage()
		return nil
	}
	client := transit.NewClient(os.Getenv("TRANSIT_API_BASE_URL"))
	switch args[0] {
	case "suggest":
		fs := flag.NewFlagSet("suggest", flag.ExitOnError)
		limit := fs.Int("limit", 10, "maximum stations to show")
		jsonOut := fs.Bool("json", false, "print raw JSON")
		_ = fs.Parse(reorderFlags(args[1:], map[string]bool{"json": true}))
		if fs.NArg() != 1 {
			return fmt.Errorf("usage: transit suggest <query>")
		}
		res, err := client.Suggest(ctx, fs.Arg(0), *limit)
		if err != nil {
			return err
		}
		if *jsonOut {
			return printJSON(res)
		}
		for _, st := range res.Stations {
			fmt.Printf("%s	%s	%s\n", st.Name, st.NameKana, st.FeedName)
		}
	case "plan":
		fs := flag.NewFlagSet("plan", flag.ExitOnError)
		date := fs.String("date", "", "service date YYYYMMDD")
		timeValue := fs.String("time", "", "HH:MM or HH:MM:SS")
		typeValue := fs.String("type", "departure", "departure|arrival|first|last")
		n := fs.Int("num", 3, "number of itineraries")
		jsonOut := fs.Bool("json", false, "print raw JSON")
		_ = fs.Parse(reorderFlags(args[1:], map[string]bool{"json": true}))
		if fs.NArg() != 2 {
			return fmt.Errorf("usage: transit plan <from> <to>")
		}
		fromID, _, err := client.ResolveStation(ctx, fs.Arg(0))
		if err != nil {
			return err
		}
		toID, _, err := client.ResolveStation(ctx, fs.Arg(1))
		if err != nil {
			return err
		}
		res, err := client.Plan(ctx, fromID, toID, transit.PlanOptions{Date: *date, Time: *timeValue, Type: *typeValue, NumItineraries: *n})
		if err != nil {
			return err
		}
		if *jsonOut {
			return printJSON(res)
		}
		printPlan(res)
	case "departures":
		fs := flag.NewFlagSet("departures", flag.ExitOnError)
		date := fs.String("date", "", "service date YYYYMMDD")
		timeValue := fs.String("time", "", "HH:MM or HH:MM:SS")
		limit := fs.Int("limit", 20, "maximum departures")
		jsonOut := fs.Bool("json", false, "print raw JSON")
		_ = fs.Parse(reorderFlags(args[1:], map[string]bool{"json": true}))
		if fs.NArg() != 1 {
			return fmt.Errorf("usage: transit departures <station>")
		}
		stationID, _, err := client.ResolveStation(ctx, fs.Arg(0))
		if err != nil {
			return err
		}
		res, err := client.Departures(ctx, stationID, transit.DepartureOptions{Date: *date, Time: *timeValue, Limit: *limit})
		if err != nil {
			return err
		}
		if *jsonOut {
			return printJSON(res)
		}
		for _, d := range res.Departures {
			headsign := d.Headsign
			if headsign != "" {
				headsign = " -> " + headsign
			}
			fmt.Printf("%s	%s	%s%s\n", transit.FormatServiceSeconds(d.DepartureSecs), transit.ModeLabel(d.Mode), d.RouteName, headsign)
		}
	case "mcp":
		return fmt.Errorf("MCP server is planned; see GitHub issues for design")
	case "version", "-v", "--version":
		fmt.Println("transit", version)
	case "help", "-h", "--help":
		usage()
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
	return nil
}

func printPlan(res transit.PlanResponse) {
	fmt.Printf("%s -> %s (%s, %s)\n", res.From.Name, res.To.Name, res.Type, res.Timezone)
	for i, j := range res.Journeys {
		fmt.Printf("\n#%d %s-%s  %d分  乗換%d回\n", i+1, transit.FormatServiceSeconds(j.DepartureSecs), transit.FormatServiceSeconds(j.ArrivalSecs), j.DurationSecs/60, j.TransferCount)
		for _, leg := range j.Legs {
			if leg.Kind == "walk" {
				fmt.Printf("  %s-%s  徒歩  %s -> %s\n", transit.FormatServiceSeconds(leg.DepartureSecs), transit.FormatServiceSeconds(leg.ArrivalSecs), leg.From.Name, leg.To.Name)
				continue
			}
			detail := strings.TrimSpace(leg.RouteName + " " + leg.Headsign)
			fmt.Printf("  %s-%s  %s  %s -> %s\n", transit.FormatServiceSeconds(leg.DepartureSecs), transit.FormatServiceSeconds(leg.ArrivalSecs), detail, leg.From.Name, leg.To.Name)
		}
	}
}

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func reorderFlags(args []string, boolFlags map[string]bool) []string {
	flags := make([]string, 0, len(args))
	positional := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
			name := strings.TrimLeft(arg, "-")
			if idx := strings.IndexByte(name, '='); idx >= 0 {
				name = name[:idx]
			}
			if !strings.Contains(arg, "=") && !boolFlags[name] && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				flags = append(flags, args[i+1])
				i++
			}
			continue
		}
		positional = append(positional, arg)
	}
	return append(flags, positional...)
}

func usage() {
	fmt.Println(`transit - Japan transit CLI powered by api.transit.ls8h.com

Usage:
  transit suggest <query> [--limit N] [--json]
  transit plan <from> <to> [--date YYYYMMDD] [--time HH:MM] [--type departure|arrival|first|last] [--json]
  transit departures <station> [--date YYYYMMDD] [--time HH:MM] [--limit N] [--json]
  transit mcp
  transit version`)
}
