package bits

// BLOCKSIZE is the size in bits of the type used as an underlying container for our BitVec.
const BLOCKSIZE = 64

// BitVec represent a vector of bits.
// bits are represented from right to left ( the low bit of the first uint64 is the first bit of the BitVec )
type BitVec = []uint64
