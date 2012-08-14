package thrift

import (
	"testing"
)

func TestCamelCase(t *testing.T) {
	cases := map[string]string{
		"test":            "Test",
		"Foo":             "Foo",
		"foo_bar":         "FooBar",
		"FooBar":          "FooBar",
		"test__ing":       "TestIng",
		"three_part_word": "ThreePartWord",
		"FOOBAR":          "FOOBAR",
		"TESTing":         "TESTing",
	}
	for k, v := range cases {
		if camelCase(k) != v {
			t.Fatalf("%s did not properly camelCase: %s", k, camelCase(k))
		}
	}
}

func BenchmarkCamelCase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		camelCase("foo_bar")
	}
}
