package storage

import (
	"encoding/binary"

	"zs-project.org/aeroview/datastructures/bits"
)

type RestartString struct {
	data      []byte
	bitVector bits.Vector
}

func (r *RestartString) ToRawStrings() []string {
	idx := 0
	var startingPoint uint64
	var l uint64
	var n int
	var n2 int

	rawStrings := make([]string, 0)
	var prev string

	for idx < len(r.data) {
		l, n = binary.Uvarint(r.data[idx:])
		startingPoint, n2 = binary.Uvarint(r.data[idx+n:])

		curr := string(r.data[idx+n+n2 : idx+n+n2+int(l)])
		rawStrings = append(rawStrings, prev[:startingPoint]+curr)
		prev = rawStrings[len(rawStrings)-1]
		idx += n + n2 + int(l)
	}

	return rawStrings
}

func FromSortedSlice(sortedStrings []string) *RestartString {
	restartPoints := make([]uint64, len(sortedStrings))
	lens := make([]uint64, len(sortedStrings))
	encodedStrings := make([]string, len(sortedStrings))

	// 1. build intermediate form
	if len(sortedStrings) != 0 {
		lens[0] = uint64(len(sortedStrings[0]))
		restartPoints[0] = 0
		encodedStrings[0] = sortedStrings[0]
	}

	for i := 1; i < len(sortedStrings); i++ {

		var smallerLen uint64
		if len(sortedStrings[i-1]) < len(sortedStrings[i]) {
			smallerLen = uint64(len(sortedStrings[i-1]))
		} else {
			smallerLen = uint64(len(sortedStrings[i]))
		}

		var restartPoint uint64
		for restartPoint < smallerLen && sortedStrings[i-1][restartPoint] == sortedStrings[i][restartPoint] {
			restartPoint++
		}
		encodedStrings[i] = sortedStrings[i][restartPoint:]
		lens[i] = uint64(len(encodedStrings[i]))
		restartPoints[i] = restartPoint
	}

	// 2. pack everything together.
	idx := 0
	bitVec := bits.NewVector(10)
	data := make([]byte, 0)
	for i := 0; i < len(sortedStrings); i++ {
		buff := make([]byte, 8)
		n := binary.PutUvarint(buff, lens[i])
		n2 := binary.PutUvarint(buff[n:], restartPoints[i])

		bitVec.Set(idx, uint8(1))
		d := append(buff[:n+n2], []byte(encodedStrings[i])...)
		data = append(data, d...)
		idx += len(d)
	}

	return &RestartString{
		data:      data,
		bitVector: *bitVec,
	}
}
