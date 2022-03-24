//go:build appengine || nounsafe

package h1

func stringToBytes(s string) []byte {
	return []byte(s)
}

func bytesToString(b []byte) string {
	return string(b)
}
