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
)

const (
	rawLength     = 12
	encodedLength = 20
	charset       = "0123456789abcdefghijklmnopqrstuv"
)

var (
	pid = os.Getpid()
	mid = getMid()

	source = NewSource()

	dec [256]byte
)

var nilID ID

var (
	ErrInvalidID = errors.New("invalid ID")
)

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
	var id ID

	i := atomic.AddUint32(&s.counter, 1)

	id[9] = byte(i >> 16)
	id[10] = byte(i >> 8)
	id[11] = byte(i)

	id[0] = mid ^ id[9]

	binary.BigEndian.PutUint64(id[1:], uint64(time.Now().UnixNano()))

	id[1], id[6], id[2], id[5] = id[6], id[1], id[5], id[2]

	id[7] = byte(pid>>8) ^ id[10]
	id[8] = byte(pid) ^ id[11]

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
	if len(text) != encodedLength {
		return ErrInvalidID
	}
	for _, c := range text {
		if dec[c] == 0xFF {
			return ErrInvalidID
		}
	}
	decode(id, text)
	return nil
}

func (id ID) Counter() uint32 {
	b := id[9:12]
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

func (id ID) Mid() uint8 {
	return id[0] ^ id[9]
}

func (id ID) Pid() uint16 {
	return binary.BigEndian.Uint16([]byte{
		id[7] ^ id[10],
		id[8] ^ id[11],
	})
}

func (id ID) Time() time.Time {
	nsec := int64(binary.BigEndian.Uint64([]byte{
		id[6],
		id[5],
		id[3],
		id[4],
		id[2],
		id[1],
		0,
		0,
	}))
	return time.Unix(0, nsec)
}

func (id ID) IsNil() bool {
	return id == nilID
}

func (id ID) Value() (driver.Value, error) {
	b, err := id.MarshalText()
	return string(b), err
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
		}
	}
}
