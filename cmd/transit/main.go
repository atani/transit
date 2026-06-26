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
		fmt.Fprintln(os.Stderr, "エラー:", err)
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
		limit := fs.Int("limit", 10, "表示する駅の最大数")
		jsonOut := fs.Bool("json", false, "生の JSON を出力")
		_ = fs.Parse(reorderFlags(args[1:], map[string]bool{"json": true}))
		if fs.NArg() != 1 {
			return fmt.Errorf("使い方: transit suggest <駅名>")
		}
		res, err := client.Suggest(ctx, fs.Arg(0), *limit)
		if err != nil {
			return err
		}
		if *jsonOut {
			return printJSON(res)
		}
		for _, st := range res.Stations {
			fmt.Printf("%s	%s	%s\n", st.Name, st.NameKana, transit.FeedLabel(st.FeedName))
		}
	case "plan":
		fs := flag.NewFlagSet("plan", flag.ExitOnError)
		date := fs.String("date", "", "対象日 YYYYMMDD")
		timeValue := fs.String("time", "", "HH:MM または HH:MM:SS")
		typeValue := fs.String("type", "departure", "departure|arrival|first|last（出発/到着/始発/終電）")
		n := fs.Int("num", 3, "経路候補の数")
		jsonOut := fs.Bool("json", false, "生の JSON を出力")
		_ = fs.Parse(reorderFlags(args[1:], map[string]bool{"json": true}))
		if fs.NArg() != 2 {
			return fmt.Errorf("使い方: transit plan <出発> <到着>")
		}
		from, to, err := client.ResolveStationPair(ctx, fs.Arg(0), fs.Arg(1))
		if err != nil {
			return err
		}
		res, err := client.Plan(ctx, from.ID, to.ID, transit.PlanOptions{Date: *date, Time: *timeValue, Type: *typeValue, NumItineraries: *n})
		if err != nil {
			return err
		}
		if *jsonOut {
			return printJSON(res)
		}
		printPlan(res)
		if planIsWalkOnly(res) {
			fmt.Fprintln(os.Stderr, "※ 直通の経路が見つかりませんでした。出発・到着が別事業者（フィード）の場合に起きます。このAPIは事業者をまたぐ乗換を計算しません。同じ事業者の駅を指定すると経路が出ることがあります。")
		}
	case "departures":
		fs := flag.NewFlagSet("departures", flag.ExitOnError)
		date := fs.String("date", "", "対象日 YYYYMMDD")
		timeValue := fs.String("time", "", "HH:MM または HH:MM:SS")
		limit := fs.Int("limit", 20, "発車案内の最大件数")
		jsonOut := fs.Bool("json", false, "生の JSON を出力")
		_ = fs.Parse(reorderFlags(args[1:], map[string]bool{"json": true}))
		if fs.NArg() != 1 {
			return fmt.Errorf("使い方: transit departures <駅名>")
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
		return fmt.Errorf("MCP サーバーは計画中です。設計は GitHub の issue を参照してください")
	case "version", "-v", "--version":
		fmt.Println("transit", version)
	case "help", "-h", "--help":
		usage()
	default:
		return fmt.Errorf("不明なコマンド: %s", args[0])
	}
	return nil
}

func printPlan(res transit.PlanResponse) {
	fmt.Printf("%s -> %s (%s, %s)\n", res.From.Name, res.To.Name, transit.TypeLabel(res.Type), res.Timezone)
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

// planIsWalkOnly reports whether the plan has no real transit, i.e. every
// itinerary is walking only (or none was found). This is the signal that the
// from/to stations are not connected within a single feed.
func planIsWalkOnly(res transit.PlanResponse) bool {
	if len(res.Journeys) == 0 {
		return true
	}
	for _, j := range res.Journeys {
		for _, leg := range j.Legs {
			if leg.Kind != "walk" {
				return false
			}
		}
	}
	return true
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
	fmt.Println(`transit - 日本の経路検索・発車案内 CLI（api.transit.ls8h.com を利用）

使い方:
  transit suggest <駅名> [--limit N] [--json]
  transit plan <出発> <到着> [--date YYYYMMDD] [--time HH:MM] [--type departure|arrival|first|last] [--json]
  transit departures <駅名> [--date YYYYMMDD] [--time HH:MM] [--limit N] [--json]
  transit mcp
  transit version`)
}
