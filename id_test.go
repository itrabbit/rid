package rid

import (
	"fmt"
	"testing"
	"time"
)

func TestIDInfo(t *testing.T) {
	s := GlobalSource()
	id, now := New(), time.Now()

	if id.Counter() != s.Counter() {
		fmt.Println("[ERROR] Invalid counter")
		t.Fail()
		return
	}
	if id.Mid() != mid {
		fmt.Println("[ERROR] Invalid NID")
		t.Fail()
		return
	}
	if id.Pid() != uint16(pid) {
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
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = New()
		}
	})
}

func BenchmarkNewString(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = New().String()
		}
	})
}

func BenchmarkFromString(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = FromString("krt54gkt2cakckbs0lm0")
		}
	})
}

func BenchmarkFromNewString(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = FromString(New().String())
		}
	})
}
