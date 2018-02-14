package rid

import (
	"bytes"
	"net"
)

func getMid() uint8 {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				return calcCRC4(i.HardwareAddr)
			}
		}
	}
	return calcCRC4([]byte{0, 0, 0, 0, 0, 128})
}
