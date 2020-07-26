package bitmap

import (
	"zs-project.org/aeroview/rank"
	"zs-project.org/aeroview/selecter"
)

type BitMap struct {
	Data rank.BitVec
	r    rank.Ranker
	s    selecter.Selecter
}

func MakeFromBitQueue(data rank.BitVec) BitMap {
	r := rank.MakeRankPopCount(data)
	return BitMap{
		Data: data,
		r:    r,
		s:    selecter.BinarySearchSelect{RankPopCount: r, Data: data},
	}

}
