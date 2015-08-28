package thrift

import (
	"testing"
)

func TestBitSet(t *testing.T) {
	s := new(BitSet)
	s.Set(13)
	s.Set(45)
	s.Clear(13)
	if s.IsSet(13) {
		t.Fatalf("Failed on Set and Clear 13")
	}
	if !s.IsSet(45) {
		t.Fatalf("Failed on set 45")
	}
	if s.IsSet(30) {
		t.Fatalf("Failed on intial 30")
	}

}
