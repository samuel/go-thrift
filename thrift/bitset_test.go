package thrift

import (
	"fmt"
	"testing"
)

func TestBitSet(t *testing.T) {
	s := new(BitSet)
	s.Set(13)
	s.Set(45)
	s.Clear(13)
	fmt.Printf("s.IsSet(13) = %t; s.IsSet(45) = %t; s.IsSet(30) = %t\n",
		s.IsSet(13), s.IsSet(45), s.IsSet(30))
	// Output: s.IsSet(13) = false; s.IsSet(45) = true; s.IsSet(30) = false
}
