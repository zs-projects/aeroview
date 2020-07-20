package storage

import (
	"encoding/binary"
)

type SuffixBlocks struct {
	references []string
	blocks     []*SuffixBlock
}

func NewSuffixBlocks(sortedStrings []string, blockSize int) *SuffixBlocks {
	references := make([]string, 0)
	blocks := make([]*SuffixBlock, 0)

	if len(sortedStrings) == 0 {
		return &SuffixBlocks{
			references: references,
			blocks:     []*SuffixBlock{},
		}
	}

	references = append(references, sortedStrings[0])
	lastIndexAdded := 0
	for i := 1; i < len(sortedStrings); i++ {
		if i%blockSize == 0 {
			block := FromStrings(sortedStrings[i-blockSize:i], references[len(references)-1])
			blocks = append(blocks, block)
			references = append(references, sortedStrings[i])
			lastIndexAdded = i
		}
	}
	block := FromStrings(sortedStrings[lastIndexAdded:], references[len(references)-1])
	blocks = append(blocks, block)

	return &SuffixBlocks{
		references: references,
		blocks:     blocks,
	}
}

// Look from which block to look for in BS fashion.
func (r *SuffixBlocks) Exists(key string) bool {
	left := 0
	right := len(r.references) - 1
	for left <= right {
		mid := (left + right) / 2
		if r.references[mid] == key {
			return true
		} else if mid == len(r.references) - 1 && r.references[mid] < key {
			return r.blocks[mid].Exists(r.references[mid], key)
		} else if mid != 0 && r.references[mid] > key && r.references[mid - 1] < key {
			return r.blocks[mid - 1].Exists(r.references[mid - 1], key)
		} else if r.references[mid] < key {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return false
}

type SuffixBlock struct {
	raws         []byte // list of non-common suffixes
	prefixPoints []byte // where to start the prefix from. encoded in uvarInt
	offsets      []byte // offsets of raws encoded in uvarint delta
}

// Exists goes thru each suffix in the block to find the key.
func (r *SuffixBlock) Exists(reference, key string) bool {
	offsetIdx := 0
	leftOffsetIndex := 0

	prefixIndex := 0

	for offsetIdx < len(r.offsets) {
		// decode offsets
		delta, nOffsets := binary.Uvarint(r.offsets[offsetIdx:])
		offsetIdx += nOffsets
		rightOffsetIndex := leftOffsetIndex + int(delta)

		// decode retart point
		prefixPoint, nPrefixes := binary.Uvarint(r.prefixPoints[prefixIndex:])
		prefixIndex += nPrefixes

		// decode string
		decoded := reference[:prefixPoint] + string(r.raws[leftOffsetIndex:rightOffsetIndex])
		if decoded == key {
			return true
		}

		leftOffsetIndex = rightOffsetIndex
	}
	return false
}

func FromStrings(sortedStrings []string, reference string) *SuffixBlock {

	raws := make([]byte, 0)
	prefixPoints := make([]byte, 0)
	idxs := make([]byte, 0)

	for _, s := range sortedStrings {

		var smallerLen int
		var prefixPoint int
		if len(s) < len(reference) {
			smallerLen = len(s)
		} else {
			smallerLen = len(reference)
		}

		for prefixPoint < smallerLen && s[prefixPoint] == reference[prefixPoint] {
			prefixPoint++
		}
		bts := []byte(s[prefixPoint:])
		raws = append(raws, bts...)

		// add prefix point
		prefixPointsBuff := make([]byte, 4)
		n := binary.PutUvarint(prefixPointsBuff, uint64(prefixPoint))
		prefixPoints = append(prefixPoints, prefixPointsBuff[:n]...)

		// add index
		indexBuff := make([]byte, 4)
		n = binary.PutUvarint(indexBuff, uint64(len(bts)))
		idxs = append(idxs, indexBuff[:n]...)
	}

	return &SuffixBlock{
		raws:         raws,
		prefixPoints: prefixPoints,
		offsets:      idxs,
	}
}
