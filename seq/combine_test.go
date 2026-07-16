package seq_test

import (
	"reflect"
	"testing"

	"github.com/everettxwzhu/go-toolkit/seq"
)

func TestConcat(t *testing.T) {
	got := seq.Of(1, 2).Concat(seq.Empty[int](), seq.Of(3, 4)).Collect()
	want := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Concat().Collect() = %v, want %v", got, want)
	}
}

func TestConcatWithoutOthers(t *testing.T) {
	got := seq.Of(1, 2).Concat().Collect()
	want := []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Concat().Collect() = %v, want %v", got, want)
	}
}
