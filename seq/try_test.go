package seq_test

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/everettxwzhu/go-toolkit/seq"
)

func TestTryCollect(t *testing.T) {
	got := seq.Of("1", "2", "3").TryCollect(strconv.Atoi)
	values, err := got.Get()

	if err != nil {
		t.Fatalf("TryCollect() error = %v, want nil", err)
	}
	if want := []int{1, 2, 3}; !reflect.DeepEqual(values, want) {
		t.Fatalf("TryCollect() values = %v, want %v", values, want)
	}

	empty, err := seq.Empty[string]().TryCollect(strconv.Atoi).Get()
	if err != nil {
		t.Fatalf("empty TryCollect() error = %v, want nil", err)
	}
	if len(empty) != 0 {
		t.Fatalf("empty TryCollect() values = %v, want empty slice", empty)
	}
}

func TestTryCollectStopsOnError(t *testing.T) {
	wantErr := errors.New("failed")
	calls := 0

	got := seq.Of(1, 2, 3).TryCollect(func(value int) (string, error) {
		calls++
		if value == 2 {
			return "partial", wantErr
		}
		return strconv.Itoa(value), nil
	})

	values, err := got.Get()
	if err != wantErr {
		t.Fatalf("TryCollect() error = %v, want %v", err, wantErr)
	}
	if values != nil {
		t.Fatalf("TryCollect() values = %v, want nil after error", values)
	}
	if calls != 2 {
		t.Fatalf("transform called %d times, want 2", calls)
	}
}

func TestTryReduce(t *testing.T) {
	got := seq.Of(1, 2, 3).TryReduce("0", func(acc string, value int) (string, error) {
		return acc + strconv.Itoa(value), nil
	})

	value, err := got.Get()
	if err != nil {
		t.Fatalf("TryReduce() error = %v, want nil", err)
	}
	if value != "0123" {
		t.Fatalf("TryReduce() value = %q, want %q", value, "0123")
	}

	empty, err := seq.Empty[int]().TryReduce("initial", func(acc string, value int) (string, error) {
		return acc + strconv.Itoa(value), nil
	}).Get()
	if err != nil {
		t.Fatalf("empty TryReduce() error = %v, want nil", err)
	}
	if empty != "initial" {
		t.Fatalf("empty TryReduce() value = %q, want %q", empty, "initial")
	}
}

func TestTryReduceStopsOnError(t *testing.T) {
	wantErr := errors.New("failed")
	calls := 0

	got := seq.Of(1, 2, 3).TryReduce(10, func(acc, value int) (int, error) {
		calls++
		if value == 2 {
			return acc, wantErr
		}
		return acc + value, nil
	})

	value, err := got.Get()
	if err != wantErr {
		t.Fatalf("TryReduce() error = %v, want %v", err, wantErr)
	}
	if value != 0 {
		t.Fatalf("TryReduce() value = %d, want zero value after error", value)
	}
	if calls != 2 {
		t.Fatalf("reducer called %d times, want 2", calls)
	}
}

func TestTryForEach(t *testing.T) {
	var visited []int
	err := seq.Of(1, 2, 3).TryForEach(func(value int) error {
		visited = append(visited, value)
		return nil
	})

	if err != nil {
		t.Fatalf("TryForEach() error = %v, want nil", err)
	}
	if want := []int{1, 2, 3}; !reflect.DeepEqual(visited, want) {
		t.Fatalf("visited = %v, want %v", visited, want)
	}
}

func TestTryForEachStopsOnError(t *testing.T) {
	wantErr := errors.New("failed")
	var visited []int

	err := seq.Of(1, 2, 3).TryForEach(func(value int) error {
		visited = append(visited, value)
		if value == 2 {
			return wantErr
		}
		return nil
	})

	if err != wantErr {
		t.Fatalf("TryForEach() error = %v, want %v", err, wantErr)
	}
	if want := []int{1, 2}; !reflect.DeepEqual(visited, want) {
		t.Fatalf("visited = %v, want %v", visited, want)
	}
}

func TestTryAny(t *testing.T) {
	calls := 0
	got := seq.Of(1, 2, 3).TryAny(func(value int) (bool, error) {
		calls++
		return value == 2, nil
	})

	value, err := got.Get()
	if err != nil || !value {
		t.Fatalf("TryAny() = (%t, %v), want (true, nil)", value, err)
	}
	if calls != 2 {
		t.Fatalf("predicate called %d times, want 2", calls)
	}

	value, err = seq.Empty[int]().TryAny(func(int) (bool, error) {
		return true, nil
	}).Get()
	if err != nil || value {
		t.Fatalf("empty TryAny() = (%t, %v), want (false, nil)", value, err)
	}
}

func TestTryAnyStopsOnError(t *testing.T) {
	wantErr := errors.New("failed")
	calls := 0

	got := seq.Of(1, 2, 3).TryAny(func(value int) (bool, error) {
		calls++
		if value == 2 {
			return false, wantErr
		}
		return false, nil
	})

	value, err := got.Get()
	if err != wantErr {
		t.Fatalf("TryAny() error = %v, want %v", err, wantErr)
	}
	if value {
		t.Fatal("TryAny() value = true, want false after error")
	}
	if calls != 2 {
		t.Fatalf("predicate called %d times, want 2", calls)
	}
}

func TestTryEvery(t *testing.T) {
	calls := 0
	got := seq.Of(2, 4, 5, 6).TryEvery(func(value int) (bool, error) {
		calls++
		return value%2 == 0, nil
	})

	value, err := got.Get()
	if err != nil || value {
		t.Fatalf("TryEvery() = (%t, %v), want (false, nil)", value, err)
	}
	if calls != 3 {
		t.Fatalf("predicate called %d times, want 3", calls)
	}

	value, err = seq.Empty[int]().TryEvery(func(int) (bool, error) {
		return false, nil
	}).Get()
	if err != nil || !value {
		t.Fatalf("empty TryEvery() = (%t, %v), want (true, nil)", value, err)
	}
}

func TestTryEveryStopsOnError(t *testing.T) {
	wantErr := errors.New("failed")
	calls := 0

	got := seq.Of(2, 4, 6).TryEvery(func(value int) (bool, error) {
		calls++
		if value == 4 {
			return false, wantErr
		}
		return true, nil
	})

	value, err := got.Get()
	if err != wantErr {
		t.Fatalf("TryEvery() error = %v, want %v", err, wantErr)
	}
	if value {
		t.Fatal("TryEvery() value = true, want false after error")
	}
	if calls != 2 {
		t.Fatalf("predicate called %d times, want 2", calls)
	}
}

func TestTryFind(t *testing.T) {
	calls := 0
	got := seq.Of(1, 2, 3).TryFind(func(value int) (bool, error) {
		calls++
		return value%2 == 0, nil
	})

	found, err := got.Get()
	if err != nil {
		t.Fatalf("TryFind() error = %v, want nil", err)
	}
	if value, ok := found.Get(); value != 2 || !ok {
		t.Fatalf("TryFind() value = (%d, %t), want (2, true)", value, ok)
	}
	if calls != 2 {
		t.Fatalf("predicate called %d times, want 2", calls)
	}

	missing, err := seq.Of(1, 3).TryFind(func(value int) (bool, error) {
		return value%2 == 0, nil
	}).Get()
	if err != nil {
		t.Fatalf("non-matching TryFind() error = %v, want nil", err)
	}
	if _, ok := missing.Get(); ok {
		t.Fatal("non-matching TryFind() returned Some, want None")
	}
}

func TestTryFindStopsOnError(t *testing.T) {
	wantErr := errors.New("failed")
	calls := 0

	got := seq.Of(1, 2, 3).TryFind(func(value int) (bool, error) {
		calls++
		if value == 2 {
			return false, wantErr
		}
		return false, nil
	})

	found, err := got.Get()
	if err != wantErr {
		t.Fatalf("TryFind() error = %v, want %v", err, wantErr)
	}
	if _, ok := found.Get(); ok {
		t.Fatal("TryFind() returned Some after error, want zero Option")
	}
	if calls != 2 {
		t.Fatalf("predicate called %d times, want 2", calls)
	}
}

func TestMapResult(t *testing.T) {
	wantErr := errors.New("failed")
	results := seq.Of(1, 2, 3).MapResult(func(value int) (string, error) {
		if value == 2 {
			return "partial", wantErr
		}
		return strconv.Itoa(value), nil
	}).Collect()

	if len(results) != 3 {
		t.Fatalf("len(MapResult().Collect()) = %d, want 3", len(results))
	}

	if value, err := results[0].Get(); value != "1" || err != nil {
		t.Fatalf("results[0] = (%q, %v), want (\"1\", nil)", value, err)
	}
	if value, err := results[1].Get(); value != "partial" || err != wantErr {
		t.Fatalf("results[1] = (%q, %v), want (\"partial\", %v)", value, err, wantErr)
	}
	if value, err := results[2].Get(); value != "3" || err != nil {
		t.Fatalf("results[2] = (%q, %v), want (\"3\", nil)", value, err)
	}
}

func TestMapResultIsLazyAndSupportsEarlyTermination(t *testing.T) {
	calls := 0
	mapped := seq.Of(1, 2, 3).MapResult(func(value int) (int, error) {
		calls++
		return value * 2, nil
	})

	if calls != 0 {
		t.Fatalf("transform called %d times before iteration, want 0", calls)
	}

	results := mapped.Take(2).Collect()
	if len(results) != 2 {
		t.Fatalf("len(MapResult().Take(2).Collect()) = %d, want 2", len(results))
	}
	if calls != 2 {
		t.Fatalf("transform called %d times after Take(2), want 2", calls)
	}
}
