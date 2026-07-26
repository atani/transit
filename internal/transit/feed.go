package transit

import (
	"regexp"
	"strings"
)

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
	"TokyuBus":      "東急バス",
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

// shardSuffix matches the feed shard markers the upstream API appends when it
// splits one operator across several feeds, e.g. "分割01". Those numbers are an
// internal detail, so they are dropped from the displayed name.
var shardSuffix = regexp.MustCompile(`^分割\d+$`)

// FeedLabel returns a human-friendly feed name. Some feeds already expose a
// clean line name (e.g. "埼京線"); others expose a raw odpt operator identifier
// such as "odpt.Operator:TokyoMetro", which this maps to a Japanese name (or,
// for unknown operators, strips the prefix). The identifier may carry a suffix
// after the operator: a shard marker like "分割01" is dropped, while a
// meaningful one like "千代田線" is kept.
func FeedLabel(feed string) string {
	const opPrefix = "odpt.Operator:"
	if !strings.HasPrefix(feed, opPrefix) {
		return feed
	}
	op, suffix, _ := strings.Cut(strings.TrimPrefix(feed, opPrefix), " ")
	name := op
	if ja, ok := odptOperators[op]; ok {
		name = ja
	}
	suffix = strings.TrimSpace(suffix)
	if suffix == "" || shardSuffix.MatchString(suffix) {
		return name
	}
	return name + " " + suffix
}
