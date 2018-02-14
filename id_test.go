package rid

import (
	"fmt"
	"testing"
	"time"
)

func TestIDInfo(t *testing.T) {
	s := GlobalSource()

	now := time.Now()
	id := New()

	fmt.Println(id)

	if id.Counter() != s.Counter() {
		fmt.Println("Invalid counter")
		t.Fail()
		return
	}
	if id.Mid() != mid {
		fmt.Println("Invalid NID")
		t.Fail()
		return
	}
	if id.Pid() != uint16(pid) {
		fmt.Println("Invalid PID")
		t.Fail()
		return
	}
	if id.Time().Unix() != now.Unix() {
		fmt.Println("Invalid Timestamp")
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
			fmt.Println(id[:])
			fmt.Println(cid[:])
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
