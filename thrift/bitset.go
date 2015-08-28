package thrift

import "errors"

const size = 64

type bits uint64

// BitSet is a set of bits that can be set, cleared and queried.
type BitSet []bits

func numChange(t interface{}) (uint64, error) {

	switch t := t.(type) {
	case int:
		return uint64(t), nil
	case uint:
		return uint64(t), nil
	case int8:
		return uint64(t), nil
	case int16:
		return uint64(t), nil
	case int32:
		return uint64(t), nil
	case int64:
		return uint64(t), nil
	case uint8:
		return uint64(t), nil
	case uint16:
		return uint64(t), nil
	case uint32:
		return uint64(t), nil
	case uint64:
		return uint64(t), nil
	}
	return 0, errors.New("type error")
}

// Set ensures that the given bit is set in the BitSet.
func (s *BitSet) Set(n interface{}) {
	i, _ := numChange(n)
	if len(*s) < int(i/size+1) {
		r := make([]bits, i/size+1)
		copy(r, *s)
		*s = r
	}
	(*s)[i/size] |= 1 << (i % size)
}

// Clear ensures that the given bit is cleared (not set) in the BitSet.
func (s *BitSet) Clear(n interface{}) {
	i, _ := numChange(n)
	if len(*s) >= int(i/size+1) {
		(*s)[i/size] &^= 1 << (i % size)
	}
}

// IsSet returns true if the given bit is set, false if it is cleared.
func (s *BitSet) IsSet(n interface{}) bool {
	i, _ := numChange(n)
	return (*s)[i/size]&(1<<(i%size)) != 0
}
func (s *BitSet) IsEmpty() bool {
	for _, val := range *s {
		if val != 0 {
			return false
		}
	}
	return true
}
