package thrift

type bitset []uint64

func bitsetCapacity(n int) int {
	return n/64 + 1
}

func newBitset(size int) *bitset {
	b := make(bitset, bitsetCapacity(size))
	return &b
}

func (b *bitset) Set(bit int) {
	for len(*b) < bitsetCapacity(bit) {
		*b = append(*b, 0)
	}
	(*b)[bit/64] |= 1 << uint64(bit%64)
}

func (b *bitset) Get(bit int) bool {
	if len(*b) < bitsetCapacity(bit) {
		return false
	}
	return (*b)[bit/64]&(1<<uint64(bit%64)) != 0
}

func (b *bitset) Clone() *bitset {
	out := make(bitset, len(*b))
	for i, bits := range *b {
		out[i] = bits
	}
	return &out
}

func (b *bitset) Clear(bit int) {
	for len(*b) < bitsetCapacity(bit) {
		*b = append(*b, 0)
	}
	(*b)[bit/64] &= ^(1 << uint64(bit%64))
}

func (b *bitset) Empty() bool {
	for _, bits := range *b {
		if bits != 0 {
			return false
		}
	}
	return true
}

// Bits returns a slice of bits that are set in the bitset.
func (b *bitset) Bits() []int {
	out := []int{}
	for i, bits := range *b {
		for j := 0; j < 64; j++ {
			if bits&(1<<uint64(j)) != 0 {
				out = append(out, i*64+int(j))
			}
		}
	}
	return out
}
