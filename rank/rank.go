package rank

// Ranker ...
type Ranker interface {
	Rank(int) int
}

// Selecter ...
type Selecter interface {
	Select(i int) int
}
