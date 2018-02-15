package rid

import (
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"sync/atomic"
	"time"
	"bytes"
	"strconv"
)

const (
	rawLength     = 12
	encodedLength = 20
	numeralLength = rawLength * 3
	charset       = "0123456789abcdefghikmnopqrstwxyz"
)

var (
	pid = os.Getpid()
	mid = getMid()

	source = NewSource()

	epoch = getEpoch()

	dec [256]byte
)

var nilID ID

var (
	ErrInvalidID = errors.New("invalid ID")
)

// 4-byte value representing the seconds since the Unix epoch
// 1-byte hardware address CRC4 ID
// 2-byte process id
// 3-byte counter
// 2-byte nanoseconds (the first 2 bytes of the four byte value)
type ID [rawLength]byte

type Source struct {
	counter uint32
}

func (s Source) Counter() uint32 {
	return s.counter
}

func (s *Source) Seed(pos uint32) {
	atomic.StoreUint32(&s.counter, pos)
}

func (s *Source) NewID() ID {

	i, id, ns := atomic.AddUint32(&s.counter, 1), ID{}, time.Now().Sub(epoch).Nanoseconds()

	id[0] = byte(ns >> 56)
	id[1] = byte(ns >> 48)
	id[2] = byte(ns >> 40)
	id[3] = byte(ns >> 32)
	id[4] = mid
	id[5] = byte(pid >> 8)
	id[6] = byte(pid)
	id[7] = byte(i >> 16)
	id[8] = byte(i >> 8)
	id[9] = byte(i)
	id[10] = byte(ns >> 24)
	id[11] = byte(ns >> 16)

	return id
}

func (id ID) String() string {
	text := make([]byte, encodedLength)
	encode(text, id[:])
	return string(text)
}

func (id ID) MarshalText() ([]byte, error) {
	text := make([]byte, encodedLength)
	encode(text, id[:])
	return text, nil
}

func (id *ID) UnmarshalText(text []byte) error {
	if len(text) == encodedLength {
		for _, c := range text {
			if dec[c] == 0xFF {
				return ErrInvalidID
			}
		}
		decode(id, text)
		return nil
	} else if len(text) == numeralLength {
		return decodeNumeral(id, text)
	} else if len(text) == rawLength {
		_, _ = id[11], text[11]
		for i, b := range text {
			id[i] = b
		}
	}
	return ErrInvalidID
}

func (id ID) NumeralString() string {
	buf := new(bytes.Buffer)
	for _,b := range id[:] {
		if b < 100 {
			buf.WriteByte('0')
		}
		if b < 10 {
			buf.WriteByte('0')
		}
		buf.WriteString(strconv.FormatUint(uint64(b), 10))
	}
	return buf.String()
}


func (id ID) Counter() uint32 {
	b := id[7:10]
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

func (id ID) Mid() uint8 {
	return id[4]
}

func (id ID) Pid() uint16 {
	return binary.BigEndian.Uint16(id[5:7])
}

func (id ID) Time() time.Time {
	ns := uint64(0) | uint64(0)<<8 | uint64(id[11])<<16 | uint64(id[10])<<24 |
		  uint64(id[3])<<32 | uint64(id[2])<<40 | uint64(id[1])<<48 | uint64(id[0])<<56

	return getEpoch().Add(time.Duration(ns))
}

func (id ID) IsNil() bool {
	return id == nilID
}

func (id ID) Value() (driver.Value, error) {
	b, err := id.MarshalText()
	return string(b), err
}

func (ID) SqlType(dialect string, size int, settings map[string]string) string {
	switch dialect {
	case "mysql":
		if _, ok := settings["NOT NULL"]; ok {
			return "VARCHAR(20)"
		}
		return "VARCHAR(20) NULL"
	default:
		return "VARCHAR(20)"
	}
}

func (id *ID) Scan(value interface{}) (err error) {
	switch val := value.(type) {
	case string:
		return id.UnmarshalText([]byte(val))
	case []byte:
		return id.UnmarshalText(val)
	default:
		return fmt.Errorf("scanning unsupported type: %T", value)
	}
}

func FromString(id string) (ID, error) {
	i := &ID{}
	err := i.UnmarshalText([]byte(id))
	return *i, err
}

func Mid() uint8 {
	return mid
}

func Pid() uint16 {
	return uint16(pid)
}

func GlobalSource() *Source {
	return source
}

func New() ID {
	return source.NewID()
}

func NewSource() *Source {
	return &Source{
		counter: randUint32(),
	}
}

// From https://github.com/rs/xid
func encode(dst, id []byte) {
	dst[0] = charset[id[0]>>3]
	dst[1] = charset[(id[1]>>6)&0x1F|(id[0]<<2)&0x1F]
	dst[2] = charset[(id[1]>>1)&0x1F]
	dst[3] = charset[(id[2]>>4)&0x1F|(id[1]<<4)&0x1F]
	dst[4] = charset[id[3]>>7|(id[2]<<1)&0x1F]
	dst[5] = charset[(id[3]>>2)&0x1F]
	dst[6] = charset[id[4]>>5|(id[3]<<3)&0x1F]
	dst[7] = charset[id[4]&0x1F]
	dst[8] = charset[id[5]>>3]
	dst[9] = charset[(id[6]>>6)&0x1F|(id[5]<<2)&0x1F]
	dst[10] = charset[(id[6]>>1)&0x1F]
	dst[11] = charset[(id[7]>>4)&0x1F|(id[6]<<4)&0x1F]
	dst[12] = charset[id[8]>>7|(id[7]<<1)&0x1F]
	dst[13] = charset[(id[8]>>2)&0x1F]
	dst[14] = charset[(id[9]>>5)|(id[8]<<3)&0x1F]
	dst[15] = charset[id[9]&0x1F]
	dst[16] = charset[id[10]>>3]
	dst[17] = charset[(id[11]>>6)&0x1F|(id[10]<<2)&0x1F]
	dst[18] = charset[(id[11]>>1)&0x1F]
	dst[19] = charset[(id[11]<<4)&0x1F]
}

// From https://github.com/rs/xid
func decode(id *ID, src []byte) {
	id[0] = dec[src[0]]<<3 | dec[src[1]]>>2
	id[1] = dec[src[1]]<<6 | dec[src[2]]<<1 | dec[src[3]]>>4
	id[2] = dec[src[3]]<<4 | dec[src[4]]>>1
	id[3] = dec[src[4]]<<7 | dec[src[5]]<<2 | dec[src[6]]>>3
	id[4] = dec[src[6]]<<5 | dec[src[7]]
	id[5] = dec[src[8]]<<3 | dec[src[9]]>>2
	id[6] = dec[src[9]]<<6 | dec[src[10]]<<1 | dec[src[11]]>>4
	id[7] = dec[src[11]]<<4 | dec[src[12]]>>1
	id[8] = dec[src[12]]<<7 | dec[src[13]]<<2 | dec[src[14]]>>3
	id[9] = dec[src[14]]<<5 | dec[src[15]]
	id[10] = dec[src[16]]<<3 | dec[src[17]]>>2
	id[11] = dec[src[17]]<<6 | dec[src[18]]<<1 | dec[src[19]]>>4
}

func decodeNumeral(id *ID, src []byte) error {
	if len(src) % 3 != 0 {
		return ErrInvalidID
	}
	for i, pos := 0, 0; i < len(src); i, pos = i + 3, pos + 1 {
		b, err := strconv.ParseUint(string(src[i:i+3]), 10, 8)
		if err != nil {
			return err
		}
		id[pos] = uint8(b)
	}
	return nil
}

func getEpoch() time.Time {
	return time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
}

func init() {
	for i := 0; i < len(dec); i++ {
		dec[i] = 0xFF
	}
	for i := 0; i < len(charset); i++ {
		dec[charset[i]] = byte(i)
	}
	if pid == 1 {
		if b, err := ioutil.ReadFile("/proc/1/cpuset"); err == nil && len(b) > 1 {
			pid = int(crc32.ChecksumIEEE(b))
		} else {
			pid = int(randUint32())
		}
	}
}
