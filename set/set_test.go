package set_test

import (
	"strconv"
	"testing"

	"github.com/everettxwzhu/go-toolkit/set"
)

func assertSet[T comparable](t *testing.T, got set.Set[T], want ...T) {
	t.Helper()

	if got.Len() != len(want) {
		t.Fatalf("Set.Len() = %d, want %d", got.Len(), len(want))
	}

	for _, value := range want {
		if !got.Contains(value) {
			t.Fatalf("Set does not contain %v", value)
		}
	}
}

func TestNew(t *testing.T) {
	s := set.New(1, 2, 2, 3)

	assertSet(t, s, 1, 2, 3)
	if s.Contains(4) {
		t.Fatal("Set.Contains(4) = true, want false")
	}

	empty := set.New[int]()
	assertSet(t, empty)
}

func TestAddAndRemove(t *testing.T) {
	s := set.New[int]()

	s.Add(1, 2, 2, 3)
	assertSet(t, s, 1, 2, 3)

	s.Remove(2, 4)
	assertSet(t, s, 1, 3)
}

func TestSetCopiesShareStorage(t *testing.T) {
	original := set.New(1, 2)
	copied := original

	copied.Add(3)
	copied.Remove(1)

	assertSet(t, original, 2, 3)
	assertSet(t, copied, 2, 3)
}

func TestClone(t *testing.T) {
	original := set.New(1, 2)
	cloned := original.Clone()

	cloned.Add(3)
	cloned.Remove(1)

	assertSet(t, original, 1, 2)
	assertSet(t, cloned, 2, 3)
}

func TestAll(t *testing.T) {
	s := set.New(1, 2, 3)
	visited := make(map[int]bool)

	for value := range s.All() {
		visited[value] = true
	}

	if len(visited) != 3 {
		t.Fatalf("All() yielded %d distinct values, want 3", len(visited))
	}
	for _, value := range []int{1, 2, 3} {
		if !visited[value] {
			t.Fatalf("All() did not yield %d", value)
		}
	}
}

func TestAllSupportsEarlyTermination(t *testing.T) {
	s := set.New(1, 2, 3)
	calls := 0

	s.All()(func(int) bool {
		calls++
		return false
	})

	if calls != 1 {
		t.Fatalf("All() yield called %d times, want 1", calls)
	}
}

func TestFilter(t *testing.T) {
	original := set.New(1, 2, 3, 4)
	filtered := original.Filter(func(value int) bool {
		return value%2 == 0
	})

	assertSet(t, filtered, 2, 4)
	assertSet(t, original, 1, 2, 3, 4)
}

func TestMap(t *testing.T) {
	original := set.New(1, 2, 3)
	mapped := original.Map(strconv.Itoa)

	assertSet(t, mapped, "1", "2", "3")
	assertSet(t, original, 1, 2, 3)

	collapsed := original.Map(func(value int) int {
		return value % 2
	})
	assertSet(t, collapsed, 0, 1)
}

func TestUnion(t *testing.T) {
	left := set.New(1, 2)
	right := set.New(2, 3)
	got := left.Union(right)

	assertSet(t, got, 1, 2, 3)
	assertSet(t, left, 1, 2)
	assertSet(t, right, 2, 3)
}

func TestIntersection(t *testing.T) {
	left := set.New(1, 2, 3)
	right := set.New(2, 3, 4, 5)
	got := left.Intersection(right)

	assertSet(t, got, 2, 3)
	assertSet(t, left, 1, 2, 3)
	assertSet(t, right, 2, 3, 4, 5)
}

func TestDifference(t *testing.T) {
	left := set.New(1, 2, 3)
	right := set.New(2, 4)
	got := left.Difference(right)

	assertSet(t, got, 1, 3)
	assertSet(t, left, 1, 2, 3)
	assertSet(t, right, 2, 4)
}

func TestSymmetricDifference(t *testing.T) {
	left := set.New(1, 2, 3)
	right := set.New(2, 3, 4)
	got := left.SymmetricDifference(right)

	assertSet(t, got, 1, 4)
	assertSet(t, left, 1, 2, 3)
	assertSet(t, right, 2, 3, 4)
}

func TestIsSubset(t *testing.T) {
	tests := []struct {
		name   string
		set    set.Set[int]
		other  set.Set[int]
		subset bool
	}{
		{name: "proper subset", set: set.New(1, 2), other: set.New(1, 2, 3), subset: true},
		{name: "equal", set: set.New(1, 2), other: set.New(1, 2), subset: true},
		{name: "missing value", set: set.New(1, 4), other: set.New(1, 2, 3), subset: false},
		{name: "larger", set: set.New(1, 2, 3), other: set.New(1, 2), subset: false},
		{name: "empty", set: set.New[int](), other: set.New(1), subset: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.IsSubset(tt.other); got != tt.subset {
				t.Fatalf("IsSubset() = %t, want %t", got, tt.subset)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name  string
		left  set.Set[int]
		right set.Set[int]
		equal bool
	}{
		{name: "same values", left: set.New(1, 2, 3), right: set.New(3, 2, 1), equal: true},
		{name: "different values", left: set.New(1, 2), right: set.New(1, 3), equal: false},
		{name: "different lengths", left: set.New(1, 2), right: set.New(1, 2, 3), equal: false},
		{name: "empty", left: set.New[int](), right: set.New[int](), equal: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.left.Equal(tt.right); got != tt.equal {
				t.Fatalf("Equal() = %t, want %t", got, tt.equal)
			}
		})
	}
}
