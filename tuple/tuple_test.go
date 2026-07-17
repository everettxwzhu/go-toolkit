package tuple_test

import (
	"strconv"
	"testing"

	"github.com/everettxwzhu/go-toolkit/tuple"
)

func TestPairZeroValue(t *testing.T) {
	var p tuple.Pair[int, string]

	if p.First != 0 || p.Second != "" {
		t.Fatalf("zero Pair = (%d, %q), want (0, \"\")", p.First, p.Second)
	}

	first, second := p.Unpack()
	if first != 0 || second != "" {
		t.Fatalf("zero Pair.Unpack() = (%d, %q), want (0, \"\")", first, second)
	}
}

func TestNewAndUnpack(t *testing.T) {
	p := tuple.New(42, "answer")

	if p.First != 42 || p.Second != "answer" {
		t.Fatalf("New() = (%d, %q), want (42, \"answer\")", p.First, p.Second)
	}

	first, second := p.Unpack()
	if first != 42 || second != "answer" {
		t.Fatalf("Unpack() = (%d, %q), want (42, \"answer\")", first, second)
	}
}

func TestSwap(t *testing.T) {
	p := tuple.New(42, "answer")
	swapped := p.Swap()

	if swapped.First != "answer" || swapped.Second != 42 {
		t.Fatalf("Swap() = (%q, %d), want (\"answer\", 42)", swapped.First, swapped.Second)
	}
	if p.First != 42 || p.Second != "answer" {
		t.Fatalf("Swap() modified original Pair to (%d, %q)", p.First, p.Second)
	}
}

func TestSwapTwiceRestoresOriginalPair(t *testing.T) {
	p := tuple.New(42, "answer")

	if got := p.Swap().Swap(); got != p {
		t.Fatalf("Swap().Swap() = %v, want %v", got, p)
	}
}

func TestMapFirst(t *testing.T) {
	p := tuple.New(42, true)
	calls := 0

	mapped := p.MapFirst(func(value int) string {
		calls++
		return strconv.Itoa(value)
	})

	if mapped.First != "42" || !mapped.Second {
		t.Fatalf("MapFirst() = (%q, %t), want (\"42\", true)", mapped.First, mapped.Second)
	}
	if calls != 1 {
		t.Fatalf("MapFirst transform called %d times, want 1", calls)
	}
	if p.First != 42 || !p.Second {
		t.Fatalf("MapFirst() modified original Pair to (%d, %t)", p.First, p.Second)
	}
}

func TestMapSecond(t *testing.T) {
	p := tuple.New("answer", 42)
	calls := 0

	mapped := p.MapSecond(func(value int) string {
		calls++
		return strconv.Itoa(value)
	})

	if mapped.First != "answer" || mapped.Second != "42" {
		t.Fatalf("MapSecond() = (%q, %q), want (\"answer\", \"42\")", mapped.First, mapped.Second)
	}
	if calls != 1 {
		t.Fatalf("MapSecond transform called %d times, want 1", calls)
	}
	if p.First != "answer" || p.Second != 42 {
		t.Fatalf("MapSecond() modified original Pair to (%q, %d)", p.First, p.Second)
	}
}
