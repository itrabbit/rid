package rid

import (
	"fmt"
	"testing"
	"time"
	"encoding/binary"
)

func TestRandomUint32(t *testing.T)  {
	if b := randUint32(); b == 0 {
		t.Fail()
		return
	}
}

func TestNanosecondsAccuracy(t *testing.T)  {

	ns := time.Now().Sub(epoch).Nanoseconds()

	var b [8]byte

	_ = b[7]
	b[0] = byte(ns >> 56)
	b[1] = byte(ns >> 48)
	b[2] = byte(ns >> 40)
	b[3] = byte(ns >> 32)
	b[4] = byte(ns >> 24)
	b[5] = byte(ns >> 16)
	b[6] = 0
	b[7] = 0

	short := int64(binary.BigEndian.Uint64(b[:]))

	if (short * 100 / ns) < 90 {
		t.Fail()
		return
	}
}

func TestID_IsNil(t *testing.T) {
	id := New()
	if id.IsNil() {
		t.Fail()
		return
	}
	id = ID{}
	if !id.IsNil() {
		t.Fail()
		return
	}
}

func TestNumeralString(t *testing.T)  {
	src := New()
	id, err := FromString(src.NumeralString())
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
		return
	}
	if id != src {
		fmt.Println(id, "!=", src)
		t.Fail()
		return
	}
}

func TestIDInfo(t *testing.T) {
	s := GlobalSource()
	id, now := New(), time.Now()

	if id.Counter() != s.Counter() {
		fmt.Println("[ERROR] Invalid counter")
		t.Fail()
		return
	}
	if id.Mid() != Mid() {
		fmt.Println("[ERROR] Invalid NID")
		t.Fail()
		return
	}
	if id.Pid() != uint16(Pid()) {
		fmt.Println("[ERROR] Invalid PID")
		t.Fail()
		return
	}
	if id.Time().Unix() != now.Unix() {
		fmt.Println("[ERROR] Invalid Timestamp")
		t.Fail()
		return
	}
}

func TestEncodeDecode(t *testing.T) {
	var id, cid ID
	for i := 0; i < 1000; i++ {
		id = New()
		cid, _ = FromString(id.String())
		if id != cid {
			fmt.Println("[ERROR]", id[:], " != ", cid[:])
			t.Fail()
			return
		}
	}
}

func BenchmarkNew(b *testing.B) {
	fmt.Println("BenchmarkNew")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = New()
		}
	})
}

func BenchmarkNewString(b *testing.B) {
	fmt.Println("BenchmarkNewString")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = New().String()
		}
	})
}

func BenchmarkNewInfo(b *testing.B) {
	fmt.Println("BenchmarkNewInfo")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := New()
			id.Mid()
			id.Pid()
			id.Time()
			id.Counter()
		}
	})
}


func BenchmarkNewNumeralString(b *testing.B) {
	fmt.Println("BenchmarkNewNumeralString")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = New().NumeralString()
		}
	})
}

func BenchmarkFromString(b *testing.B) {
	fmt.Println("BenchmarkFromString")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = FromString("2m9p3besa0bo6mtqdss0")
		}
	})
}
