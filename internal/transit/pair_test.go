package transit

import "testing"

func TestPickFeedPair(t *testing.T) {
	st := func(id, feed string) Station {
		return Station{ID: id, FeedID: feed, Kind: "station"}
	}
	stop := func(id, feed string) Station {
		return Station{ID: id, FeedID: feed, Kind: "stop"}
	}

	t.Run("prefers two rail stations sharing a feed", func(t *testing.T) {
		from := []Station{st("a1", "feedA"), st("b1", "feedB")}
		to := []Station{st("c1", "feedC"), st("b2", "feedB")}
		f, tt := pickFeedPair(from, to)
		if f.ID != "b1" || tt.ID != "b2" {
			t.Fatalf("got (%s,%s), want (b1,b2)", f.ID, tt.ID)
		}
	})

	// A bus stop may outrank the rail station in the suggestion list; picking it
	// yields an all-bus itinerary, so rail stations win even across feeds.
	t.Run("prefers rail stations over a higher-ranked bus stop", func(t *testing.T) {
		from := []Station{st("shibuya-rail", "jr")}
		to := []Station{stop("oshiage-bus", "bus"), st("oshiage-rail", "toei")}
		f, tt := pickFeedPair(from, to)
		if f.ID != "shibuya-rail" || tt.ID != "oshiage-rail" {
			t.Fatalf("got (%s,%s), want (shibuya-rail,oshiage-rail)", f.ID, tt.ID)
		}
	})

	// A shared bus feed must not beat a rail pair on different feeds: the API
	// routes across feeds, so two stations are the better choice.
	t.Run("rail across feeds beats a shared bus feed", func(t *testing.T) {
		from := []Station{st("rail-a", "railA"), stop("bus-a", "busX")}
		to := []Station{st("rail-b", "railB"), stop("bus-b", "busX")}
		f, tt := pickFeedPair(from, to)
		if f.ID != "rail-a" || tt.ID != "rail-b" {
			t.Fatalf("got (%s,%s), want (rail-a,rail-b)", f.ID, tt.ID)
		}
	})

	t.Run("falls back to a shared feed when no rail station exists", func(t *testing.T) {
		from := []Station{stop("p", "busX"), stop("p2", "busY")}
		to := []Station{stop("q", "busY")}
		f, tt := pickFeedPair(from, to)
		if f.ID != "p2" || tt.ID != "q" {
			t.Fatalf("got (%s,%s), want (p2,q)", f.ID, tt.ID)
		}
	})

	t.Run("falls back to the top candidate of each", func(t *testing.T) {
		from := []Station{stop("x", "f1")}
		to := []Station{stop("y", "f2")}
		f, tt := pickFeedPair(from, to)
		if f.ID != "x" || tt.ID != "y" {
			t.Fatalf("got (%s,%s), want (x,y)", f.ID, tt.ID)
		}
	})

	// Passthrough inputs (geo: or feed-qualified IDs) carry no kind or feed.
	t.Run("handles candidates without kind or feed", func(t *testing.T) {
		from := []Station{{ID: "geo:35.6,139.7"}}
		to := []Station{{ID: "geo:35.7,139.8"}}
		f, tt := pickFeedPair(from, to)
		if f.ID != "geo:35.6,139.7" || tt.ID != "geo:35.7,139.8" {
			t.Fatalf("got (%s,%s)", f.ID, tt.ID)
		}
	})
}
