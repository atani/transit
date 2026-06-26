package transit

// TypeLabel returns a Japanese label for a plan type. Unknown values are
// returned unchanged.
func TypeLabel(t string) string {
	switch t {
	case "departure":
		return "出発"
	case "arrival":
		return "到着"
	case "first":
		return "始発"
	case "last":
		return "終電"
	default:
		return t
	}
}
