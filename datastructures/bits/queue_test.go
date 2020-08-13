package bits

import (
	"testing"
)

func TestBitQueuePushBack(t *testing.T) {
	// Tests that push back is working as expected.
	b := MakeQueue()
	b.PushBack(1)
	vals := b.Data()
	if len(vals) != 1 {
		t.Errorf("Expected the length of the bitQueue to be one. got %v", b.Len())
	}
	if vals[0] != 1 {
		t.Errorf("Expected the first value of the bitQueue to be one. got %v", vals[0])
	}
	b.PushBack(0)
	if vals[0] != 1 {
		t.Errorf("Expected the first value of the bitQueue to be 10 got %#b", vals[0])
	}
	b.PushBack(0)
	b.PushBack(1)
	if b.data[0] != 9 {
		t.Errorf("Expected the first value of the bitQueue to be %b ans was %#b", 9, b.data[0])
	}
	b.PushBack(0)
	b.PushBack(0)
	b.PushBack(1)
	b.PushBack(0)
	if b.data[0] != 73 {
		t.Errorf("Expected the first value of the bitQueue to be 0b10010010 ans was %#b", b.data[0])
	}
	for i := 0; i < 56; i++ {
		b.PushBack(0)
	}
	b.PushBack(1)
	if len(b.Data()) != 2 {
		t.Errorf("Expected the length of the bitQueue to be two. got %v", b.Len())
	}
	if b.Data()[1] != 1 {
		t.Errorf("Expected the second value of the bitQueue to be one. got %#b", b.Data()[1])
	}
	// Tests that Pushback is always only using the least significat bit.
	b = MakeQueue()
	b.PushBack(2)
	vals = b.Data()
	if len(vals) != 1 {
		t.Errorf("Expected the length of the bitQueue to be one. got %v", b.Len())
	}
	if vals[0] != 0 {
		t.Errorf("Expected the first value of the bitQueue to be 0. got %v", vals[0])
	}
}

func TestBitQueueAppend(t *testing.T) {
	b := MakeQueue()
	b.PushBack(1)
	b.PushBack(0)
	b.PushBack(1)

	b.Append([]uint64{uint64(0b10110100)}, 6)

	if b.data[0] != uint64(0b00110100101) {
		t.Errorf("BitQueue Append Failed.")
	}
}

func TestBitQueueToggle(t *testing.T) {
	b := MakeQueue()
	b.PushBack(1)
	b.PushBack(0)
	b.PushBack(1)

	b.Append([]uint64{uint64(0b10110100)}, 6)
	b.Toggle(0)
	b.Toggle(3)
	b.Toggle(4)
	if b.data[0] != uint64(0b00110111100) {
		t.Errorf("BitQueue Set Failed.")
	}
}

func TestBitQueueHigh(t *testing.T) {
	b := MakeQueue()
	b.PushBack(1)
	b.PushBack(0)
	b.PushBack(1)

	b.Append([]uint64{uint64(0b10110100)}, 6)
	b.High(6)
	b.High(0)
	if b.data[0] != uint64(0b00111100101) {
		t.Errorf("BitQueue Set Failed.")
	}
}

func TestBitQueueGet(t *testing.T) {
	b := MakeQueue()
	b.PushBack(1)
	b.PushBack(0)
	b.PushBack(1)

	b.Append([]uint64{uint64(0b10110100)}, 6)
	if b.data[0] != uint64(0b00110100101) {
		t.Errorf("BitQueue Set Failed.")
	}
	if b.Get(0) != uint64(0b1) {
		t.Errorf("BitQueue Set Failed.")
	}
	if b.Get(6) != uint64(0b0) {
		t.Errorf("BitQueue Set Failed.")
	}

}
