package thrift

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestBitsetCapacity(t *testing.T) {
	table := []struct{ bits, expected int }{
		{0, 1},
		{8, 1},
		{63, 1},
		{64, 2},
		{65, 2},
	}
	for _, row := range table {
		if bitsetCapacity(row.bits) != row.expected {
			t.Fatalf("expected bitset capacity for %d to be %d but was %d",
				row.bits, row.expected, bitsetCapacity(row.bits))
		}
	}
}

func TestBitset(t *testing.T) {
	table := []struct{ bit, length int }{
		{0, 1},
		{63, 1},
		{64, 2},
		{128, 3},
	}
	for _, row := range table {
		b := newBitset(0)

		b.Set(row.bit)
		if len(*b) != row.length {
			t.Fatalf("expected bitset to be length %d when bit %d is set but was %d",
				row.length, row.bit, len(*b))
		}

		if b.Empty() {
			t.Fatalf("expected bitset to not be empty")
		}

		actualBits := b.Bits()
		expectedBits := []int{row.bit}
		if !reflect.DeepEqual(actualBits, expectedBits) {
			t.Fatalf("Expected bits to be:\n%s\nBut was:\n%s", pprint(expectedBits), pprint(actualBits))
		}

		b.Clear(row.bit)
		if !b.Empty() {
			t.Fatalf("expected bitset to not be empty")
		}
	}
}

func pprint(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
