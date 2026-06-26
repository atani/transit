package transit

import "strings"

// odptOperators maps raw odpt operator identifiers to Japanese names.
var odptOperators = map[string]string{
	"TokyoMetro":    "東京メトロ",
	"Toei":          "都営",
	"JR-East":       "JR東日本",
	"JR-Central":    "JR東海",
	"JR-West":       "JR西日本",
	"Keio":          "京王",
	"Odakyu":        "小田急",
	"Tokyu":         "東急",
	"Keikyu":        "京急",
	"Keisei":        "京成",
	"Seibu":         "西武",
	"Tobu":          "東武",
	"Sotetsu":       "相鉄",
	"TWR":           "りんかい線",
	"TamaMonorail":  "多摩モノレール",
	"Yurikamome":    "ゆりかもめ",
	"TokyoMonorail": "東京モノレール",
	"MIR":           "つくばエクスプレス",
}

// FeedLabel returns a human-friendly feed name. Some feeds already expose a
// clean line name (e.g. "埼京線"); others expose a raw odpt operator identifier
// such as "odpt.Operator:TokyoMetro", which this maps to a Japanese name (or,
// for unknown operators, strips the prefix).
func FeedLabel(feed string) string {
	const opPrefix = "odpt.Operator:"
	if !strings.HasPrefix(feed, opPrefix) {
		return feed
	}
	op := strings.TrimPrefix(feed, opPrefix)
	if ja, ok := odptOperators[op]; ok {
		return ja
	}
	return op
}
