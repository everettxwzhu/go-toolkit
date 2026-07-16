package seq_test

import (
	"reflect"
	"testing"

	"github.com/everettxwzhu/go-toolkit/seq"
)

func TestSeqZeroValue(t *testing.T) {
	var s seq.Seq[int]
	if got := s.Collect(); len(got) != 0 {
		t.Fatalf("zero Seq.Collect() = %v, want empty slice", got)
	}
}

func TestAllSupportsEarlyTermination(t *testing.T) {
	visited := make([]int, 0, 2)
	seq.Range(1, 5).All()(func(value int) bool {
		visited = append(visited, value)
		return len(visited) < 2
	})

	if want := []int{1, 2}; !reflect.DeepEqual(visited, want) {
		t.Fatalf("visited = %v, want %v", visited, want)
	}
}
