package bits

import (
	"fmt"
)

// Queue implements a FIFO queue on top of BitVec.
type Queue struct {
	data Vector
	// cursor represent the position of the next free bit.
	cursor int
	// remainingCapacity is the number of unused bits in the current word.
	remainingCapacity int8
}

// MakeQueue creates a BitQueue with 0 capacity
func MakeQueue() Queue {
	return Queue{
		data:              make([]uint64, 0),
		remainingCapacity: 0,
		cursor:            0,
	}
}

// MakeBitQueueFromSlice makes a BitQueue from a slice.
func MakeBitQueueFromSlice(b []uint64, size int) (Queue, error) {
	if len(b)*BLOCKSIZE-size > BLOCKSIZE-1 {
		return Queue{}, fmt.Errorf("size is invalid, it should be between %v and %v", len(b)*8, len(b)*8-7)
	}
	return Queue{
		data:              b,
		remainingCapacity: int8(len(b)*BLOCKSIZE - size),
		cursor:            0,
	}, nil
}

// Len returns the number of bits stored.
func (m Queue) Vector() Vector {
	return Vector(m.data)
}

// Len returns the number of bits stored.
func (m Queue) Len() int {
	return len(m.data)*BLOCKSIZE - int(m.remainingCapacity)
}

// PushBack add the provided bit at the end of the BitQueue
// bit should be 0 or 1 ( otherwise the function will take the lowest bit value anyway )
// Assumes little endian encoding.
func (m *Queue) PushBack(bit uint64) {
	if m.remainingCapacity > 0 {
		newValue := m.data[len(m.data)-1] | uint64((bit&0b1)<<(BLOCKSIZE-m.remainingCapacity))
		m.data[len(m.data)-1] = newValue
		m.remainingCapacity--
	} else {
		m.data = append(m.data, uint64((bit & 0b1)))
		m.remainingCapacity = (BLOCKSIZE - 1)
	}
}

// Pop pops one bit from the queue
func (m *Queue) Pop() uint64 {
	v := m.Peek()
	m.cursor++
	return v
}

// Peek return the next from the queue withouth removing it.
func (m Queue) Peek() uint64 {
	return m.Get(m.cursor)
}

// Empty return true if the queue is empty
func (m Queue) Empty() bool {
	return m.cursor >= m.Len()
}

// Data returns the a copy of the underlying data as []byte.
func (m Queue) Data() []uint64 {
	b := make([]uint64, len(m.data))
	copy(b, m.data)
	return b
}

// Reset resets the queue.
func (m *Queue) Reset() {
	m.cursor = 0
}

// Append and the bytes in the slice to the buffer.
func (m *Queue) Append(data []uint64, size int) {
	b, _ := MakeBitQueueFromSlice(data, size)
	for !b.Empty() {
		m.PushBack(b.Pop())
	}
}

// Get return the bit value at position int
func (m Queue) Get(i int) uint64 {
	position := i / BLOCKSIZE
	offset := i % BLOCKSIZE
	return m.data[position] >> offset & 0b1
}

// Toggle set the bit balue at position i
// does nothing if i is out of bound
func (m Queue) Toggle(i int) {
	position := i / BLOCKSIZE
	if position < len(m.data) {
		offset := i % BLOCKSIZE
		m.data[position] = m.data[position] ^ 1<<(offset)
	}
}

// High set the bit balue at position i
// does nothing if i is out of bound
func (m Queue) High(i int) {
	position := i / BLOCKSIZE
	if position < len(m.data) {
		offset := i % BLOCKSIZE
		m.data[position] = m.data[position] | 1<<offset
	}
}
