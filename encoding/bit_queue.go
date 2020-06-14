package encoding

import (
	"fmt"
)

// BitQueue A queue to store bits with a []]byte underlying storage.
type BitQueue struct {
	bits              []byte
	remainingCapacity uint8
	cursor            int
}

// MakeBitQueue creates a BitQueue with 0 capacity
func MakeBitQueue() BitQueue {
	return BitQueue{
		bits:              make([]byte, 0),
		remainingCapacity: 0,
		cursor:            0,
	}
}

// MakeBitQueueFromSlice makes a BitQueue from a slice.
func MakeBitQueueFromSlice(b []byte, size int) (BitQueue, error) {
	if len(b)*8-size > 7 {
		return BitQueue{}, fmt.Errorf("size is invalid, it should be between %v and %v", len(b)*8, len(b)*8-7)
	}
	return BitQueue{
		bits:              b,
		remainingCapacity: uint8(len(b)*8 - size),
		cursor:            0,
	}, nil
}

// Len returns the number of bits stored.
func (m BitQueue) Len() int {
	return len(m.bits)*8 - int(m.remainingCapacity)
}

// PushBack add the provided bit at the end of the BitQueue
// bit should be 0 or 1 ( otherwise the function will take the lowest bit value anyway )
// Assumes little endian encoding.
func (m *BitQueue) PushBack(bit uint8) {
	if m.remainingCapacity > 0 {
		newValue := m.bits[len(m.bits)-1] | ((bit & 0b1) << (m.remainingCapacity - 1))
		m.bits[len(m.bits)-1] = newValue
		m.remainingCapacity--
	} else {
		m.bits = append(m.bits, (bit&0b1)<<7)
		m.remainingCapacity = 7
	}
}

// Pop pops one bit from the queue
func (m *BitQueue) Pop() uint8 {
	v := m.Peek()
	m.cursor++
	return v
}

// Peek return the next from the queue withouth removing it.
func (m BitQueue) Peek() uint8 {
	position := m.cursor / 8
	offset := m.cursor % 8
	return m.bits[position] >> (7 - offset) & 0b1
}

// Empty return true if the queue is empty
func (m BitQueue) Empty() bool {
	return m.cursor >= m.Len()
}

// Data returns the underlying data as []byte.
func (m BitQueue) Data() []byte {
	return m.bits
}
