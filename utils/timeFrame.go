package utils

import "fmt"

func GenTimeFrame(ms int64) string {
	if ms <= 0 {
		return ""
	}

	sec := ms / 1000

	switch {
	case sec < 60:
		return fmt.Sprintf("%dS", sec)

	case sec%3600 == 0:
		return fmt.Sprintf("%dH", sec/3600)

	case sec%60 == 0:
		return fmt.Sprintf("%dMIN", sec/60)

	case sec%86400 == 0:
		return fmt.Sprintf("%dD", sec/86400)

	default:
		return fmt.Sprintf("%dMS", ms)
	}
}
