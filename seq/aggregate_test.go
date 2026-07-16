package seq_test

import (
	"reflect"
	"testing"

	"github.com/everettxwzhu/go-toolkit/seq"
)

func TestGroupBy(t *testing.T) {
	got := seq.Of("one", "two", "three", "six").GroupBy(func(value string) int {
		return len(value)
	})
	want := map[int][]string{3: {"one", "two", "six"}, 5: {"three"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GroupBy() = %v, want %v", got, want)
	}
}

func TestToMapLastValueWins(t *testing.T) {
	got := seq.Of("one", "two", "three").ToMap(func(value string) (int, string) {
		return len(value), value
	})
	want := map[int]string{3: "two", 5: "three"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ToMap() = %v, want %v", got, want)
	}
}
