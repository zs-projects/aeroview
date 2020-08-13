package datastructures

import (
	"fmt"
	"math"

	"zs-project.org/aeroview/rank"
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

// AsBitVec returns a BitVec from a bitQueue
func (m BitQueue) AsBitVec() rank.BitVec {
	l := int(math.Max(float64(len(m.bits)/8), 1))
	vec := make([]uint64, l)
	for k, v := range m.bits {
		idx := k / 8
		switch k % 8 {
		case 0:
			u := uint64(v) << (63 - 8*1)
			vec[idx] += u
		case 1:
			u := uint64(v) << (63 - 8*2)
			vec[idx] += u
		case 2:
			u := uint64(v) << (63 - 8*3)
			vec[idx] += u
		case 3:
			u := uint64(v) << (63 - 8*4)
			vec[idx] += u
		case 4:
			u := uint64(v) << (63 - 8*5)
			vec[idx] += u
		case 5:
			u := uint64(v) << (63 - 8*6)
			vec[idx] += u
		case 6:
			u := uint64(v) << (63 - 8*7)
			vec[idx] += u
		case 7:
			u := uint64(v)
			vec[idx] += u
		}
	}
	return rank.BitVec(vec)
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
	return m.Get(m.cursor)
}

// Empty return true if the queue is empty
func (m BitQueue) Empty() bool {
	return m.cursor >= m.Len()
}

// Data returns the a copy of the underlying data as []byte.
func (m BitQueue) Data() []byte {
	b := make([]byte, len(m.bits))
	copy(b, m.bits)
	return b
}

// Reset resets the queue.
func (m *BitQueue) Reset() {
	m.cursor = 0
}

// Append and the bytes in the slice to the buffer.
func (m *BitQueue) Append(data []byte, size int) {
	b, _ := MakeBitQueueFromSlice(data, size)
	for !b.Empty() {
		m.PushBack(b.Pop())
	}
}

// Get return the bit value at position int
func (m BitQueue) Get(i int) uint8 {
	position := i / 8
	offset := i % 8
	return m.bits[position] >> (7 - offset) & 0b1
}

// Toggle set the bit balue at position i
// does nothing if i is out of bound
func (m BitQueue) Toggle(i int) {
	position := i / 8
	if position < len(m.bits) {
		offset := i % 8
		m.bits[position] = m.bits[position] ^ 1<<(7-offset)
	}
}

// High set the bit balue at position i
// does nothing if i is out of bound
func (m BitQueue) High(i int) {
	position := i / 8
	if position < len(m.bits) {
		offset := i % 8
		m.bits[position] = m.bits[position] | 1<<(7-offset)
	}
}
