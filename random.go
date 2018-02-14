package rid

import (
	crypt "crypto/rand"
	"math/rand"
	"time"
)

func randUint32() uint32 {

	rand.Int()

	b := make([]byte, 3)
	if _, err := crypt.Reader.Read(b); err != nil {
		rand.Seed(rand.Int63())
		if _, err := rand.Read(b); err != nil {
			return uint32(time.Now().Nanosecond())
		}
	}
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}
