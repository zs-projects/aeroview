package encoding

import (
	"testing"
)

func TestBitQueuePushBack(t *testing.T) {
	// Tests that push back is working as expected.
	b := MakeBitQueue()
	b.PushBack(1)
	vals := b.Data()
	if len(vals) != 1 {
		t.Errorf("Expected the length of the bitQueue to be one. got %v", b.Len())
	}
	if vals[0] != 0b10000000 {
		t.Errorf("Expected the first value of the bitQueue to be one. got %v", vals[0])
	}
	b.PushBack(0)
	if vals[0] != 0b10000000 {
		t.Errorf("Expected the first value of the bitQueue to be 10 got %#b", vals[0])
	}
	b.PushBack(0)
	b.PushBack(1)
	if b.bits[0] != 0b10010000 {
		t.Errorf("Expected the first value of the bitQueue to be 0b10010000 ans was %#b", b.bits[0])
	}
	b.PushBack(0)
	b.PushBack(0)
	b.PushBack(1)
	b.PushBack(0)
	if b.bits[0] != 0b10010010 {
		t.Errorf("Expected the first value of the bitQueue to be 0b10010010 ans was %#b", b.bits[0])
	}
	b.PushBack(1)
	if len(b.Data()) != 2 {
		t.Errorf("Expected the length of the bitQueue to be two. got %v", b.Len())
	}
	if b.Data()[1] != 0b10000000 {
		t.Errorf("Expected the second value of the bitQueue to be one. got %#b", b.Data()[1])
	}
	// Tests that Pushback is always only using the least significat bit.
	b = MakeBitQueue()
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
	b := MakeBitQueue()
	b.PushBack(1)
	b.PushBack(0)
	b.PushBack(1)

	b.Append([]byte{0b10110100}, 6)

	if b.bits[0] != 0b10110110 || b.bits[1] != 0b10000000 {
		t.Errorf("BitQueue Append Failed.")
	}
}
