package transit

import "testing"

func TestPickFeedPair(t *testing.T) {
	t.Run("prefers a shared feed", func(t *testing.T) {
		from := []Station{{ID: "a1", FeedID: "feedA"}, {ID: "b1", FeedID: "feedB"}}
		to := []Station{{ID: "c1", FeedID: "feedC"}, {ID: "b2", FeedID: "feedB"}}
		f, tt := pickFeedPair(from, to)
		if f.ID != "b1" || tt.ID != "b2" {
			t.Fatalf("got (%s,%s), want (b1,b2)", f.ID, tt.ID)
		}
	})

	t.Run("no shared feed falls back to top of each", func(t *testing.T) {
		from := []Station{{ID: "x", FeedID: "f1"}}
		to := []Station{{ID: "y", FeedID: "f2"}}
		f, tt := pickFeedPair(from, to)
		if f.ID != "x" || tt.ID != "y" {
			t.Fatalf("got (%s,%s), want (x,y)", f.ID, tt.ID)
		}
	})

	t.Run("empty feed ids do not match each other", func(t *testing.T) {
		from := []Station{{ID: "p", FeedID: ""}}
		to := []Station{{ID: "q", FeedID: ""}}
		f, tt := pickFeedPair(from, to)
		if f.ID != "p" || tt.ID != "q" {
			t.Fatalf("got (%s,%s), want (p,q)", f.ID, tt.ID)
		}
	})
}
