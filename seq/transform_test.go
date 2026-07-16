package seq_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/everettxwzhu/go-toolkit/seq"
)

func TestMap(t *testing.T) {
	got := seq.Of(1, 2, 3).Map(strconv.Itoa).Collect()
	want := []string{"1", "2", "3"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Map().Collect() = %v, want %v", got, want)
	}
}

func TestFilterMap(t *testing.T) {
	got := seq.Of("1", "bad", "3").FilterMap(func(value string) (int, bool) {
		parsed, err := strconv.Atoi(value)
		return parsed, err == nil
	}).Collect()
	want := []int{1, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FilterMap().Collect() = %v, want %v", got, want)
	}
}

func TestFlatMap(t *testing.T) {
	got := seq.Of(1, 3).FlatMap(func(value int) seq.Seq[int] {
		return seq.Range(0, value)
	}).Collect()
	want := []int{0, 0, 1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FlatMap().Collect() = %v, want %v", got, want)
	}
}
