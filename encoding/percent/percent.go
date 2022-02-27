package percent

func Decode(buffer []byte) []byte {
	writeIndex := 0

	for readIndex := 0; readIndex < len(buffer); readIndex++ {
		if buffer[readIndex] == '%' {
			if readIndex+2 >= len(buffer) {
				break
			}
			buffer[writeIndex] = DecodeHexTwo(buffer[readIndex+1], buffer[readIndex+2])
			readIndex += 2
		} else {
			buffer[writeIndex] = buffer[readIndex]
		}
		writeIndex++
	}

	return buffer[:writeIndex]
}
