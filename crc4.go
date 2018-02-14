package rid

func calcCRC4(b []byte) uint8 {
	c := 0xFFFF
	for bI := 0; bI < len(b); bI++ {
		bit := uint8(0x80)
		for bitI := 0; bitI < 8; bitI++ {
			xor := (c & 0x8000) == 0x8000
			c = c << 1
			if ((b[bI] & bit) ^ uint8(0xFF)) != uint8(0xFF) {
				c = c + 1
			}
			if xor {
				c = c ^ 0x1021
			}
			bit = bit >> 1
		}
	}
	return uint8(c)
}
