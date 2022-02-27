package percent

// https://stackoverflow.com/questions/33925589/fastest-way-to-decode-a-hexadecimal-digit
// https://stackoverflow.com/a/33925971
func DecodeHexOne(d byte) byte {
	return (d & 0x0f) + (d >> 6) + ((d >> 6) << 3)
}

func DecodeHexTwo(a, b byte) byte {
	return DecodeHexOne(a)<<4 | DecodeHexOne(b)
}
