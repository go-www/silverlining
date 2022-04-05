//go:build faststricmp
// +build faststricmp

package h1

func stricmp(a, b []byte) bool {
	return string(a) == string(b)
}
