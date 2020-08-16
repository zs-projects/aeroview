package rank

// GetRanker provide a rank operation and a get operation over a bitset.
type GetRanker interface {
	// Returns the position of the i'th bit in the underlying bitset.
	Get(i int) int
	// Returns the number of high bits between [0, position] in the underlying bitset.
	Rank(position int) int
}

// GetRankSelecter provide get, rank and select operations over a bitset.
type GetRankSelecter interface {
	// Returns the position of the i'th bit in the underlying bitset.
	Get(i int) int
	// Returns the number of high bits between [0, position] in the underlying bitset.
	Rank(position int) int
	// Returns the position of the i'th high bit in the underlying bitset.
	Select(i int) int
}

// Ranker provides a rank operation over a bitset.
type Ranker interface {
	// Returns the number of high bits between [0, position] in the underlying bitset.
	Rank(int) int
}

// Selecter provides a select operation over a bitset.
type Selecter interface {
	// Returns the position of the i'th high bit in the underlying bitset.
	Select(i int) int
}
