package h1

func stricmp(a, b []byte) bool {
	// Fast Path
	if len(a) != len(b) {
		return false
	}
	if string(a) == string(b) {
		return true
	}

	// Slow Path
	for i := 0; i < len(a); i++ {
		if !( /* case-insensitive */ (a[i] | 0x20) == (b[i] | 0x20)) {
			return false
		}
	}
	return true
}
