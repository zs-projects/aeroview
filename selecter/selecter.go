package selecter

type Selecter interface {
	Select(i uint64) uint64
}
