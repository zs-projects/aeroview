package randutils

import (
	"encoding/binary"
	"math/rand"
	"sort"
)

func RandSlice32(size int) []uint32 {
	out := make([]uint32, size)
	for i := range out {
		out[i] = rand.Uint32()
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func RandSlice64(size int) []uint64 {
	out := make([]uint64, size)
	for i := range out {
		out[i] = rand.Uint64()
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func RandSliceMixed(size int) []uint64 {
	out := make([]uint64, size)
	for i := range out {
		switch rand.Uint32() % 6 {
		case 0:
			out[i] = rand.Uint64()
		case 1, 3, 4, 5:
			out[i] = uint64(rand.Uint32())
		case 2:
			out[i] = uint64(uint16(rand.Uint32()))
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func RandSlice16(size int) []uint16 {
	out := make([]uint16, size)
	for i := range out {
		out[i] = uint16(rand.Uint32())
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func VarInt16(slice []uint16) []byte {
	out := make([]byte, len(slice)*4)
	i := 0
	for _, k := range slice {
		i += binary.PutUvarint(out[i:], uint64(k))
	}
	return out[:i]
}

func VarInt32(slice []uint32) []byte {
	out := make([]byte, len(slice)*8)
	i := 0
	for _, k := range slice {
		i += binary.PutUvarint(out[i:], uint64(k))
	}
	return out[:i]
}

func VarInt64(slice []uint64) []byte {
	out := make([]byte, len(slice)*12)
	i := 0
	for _, k := range slice {
		i += binary.PutUvarint(out[i:], uint64(k))
	}
	return out[:i]
}

func Slice32To64(slice []uint32) []uint64 {
	out := make([]uint64, 0, len(slice))
	for _, v := range slice {
		out = append(out, uint64(v))
	}
	return out
}

func Slice16To64(slice []uint16) []uint64 {
	out := make([]uint64, 0, len(slice))
	for _, v := range slice {
		out = append(out, uint64(v))
	}
	return out
}
