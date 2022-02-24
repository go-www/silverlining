//go:build appengine || nounsafe

package silverlining

func stringToBytes(s string) []byte {
	return []byte(s)
}

func bytesToString(b []byte) string {
	return string(b)
}
