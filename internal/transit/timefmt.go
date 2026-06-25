package transit

import "fmt"

// FormatServiceSeconds converts seconds from service-date midnight into HH:MM.
// It preserves after-midnight and previous-day service with +Nd/-Nd suffixes.
func FormatServiceSeconds(seconds int) string {
	day := 0
	for seconds < 0 {
		seconds += 86400
		day--
	}
	for seconds >= 86400 {
		seconds -= 86400
		day++
	}
	h := seconds / 3600
	m := (seconds % 3600) / 60
	base := fmt.Sprintf("%02d:%02d", h, m)
	if day > 0 {
		return fmt.Sprintf("%s(+%dd)", base, day)
	}
	if day < 0 {
		return fmt.Sprintf("%s(%dd)", base, day)
	}
	return base
}
