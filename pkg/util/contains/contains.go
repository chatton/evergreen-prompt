package contains

func String(strings []string, s string) bool {
	for _, ss := range strings {
		if ss == s {
			return true
		}
	}
	return false
}
