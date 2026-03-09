package array

import (
	"testing"
)

func TestFilter(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5, 6}
	evens := Filter(numbers, func(n int) bool {
		return n%2 == 0
	})

	if len(evens) != 3 {
		t.Errorf("Expected 3 even numbers, got %d", len(evens))
	}

	expected := []int{2, 4, 6}
	for i, v := range evens {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestMap(t *testing.T) {
	numbers := []int{1, 2, 3, 4}
	doubled := Map(numbers, func(n int) int {
		return n * 2
	})

	expected := []int{2, 4, 6, 8}
	for i, v := range doubled {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestReduce(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	sum := Reduce(numbers, func(acc, n int) int {
		return acc + n
	}, 0)

	if sum != 15 {
		t.Errorf("Expected sum to be 15, got %d", sum)
	}
}

func TestFind(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	result, found := Find(numbers, func(n int) bool {
		return n > 3
	})

	if !found {
		t.Error("Expected to find a number")
	}

	if result != 4 {
		t.Errorf("Expected to find 4, got %d", result)
	}
}

func TestFindNotFound(t *testing.T) {
	numbers := []int{1, 2, 3}
	_, found := Find(numbers, func(n int) bool {
		return n > 10
	})

	if found {
		t.Error("Expected not to find a number")
	}
}

func TestSome(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	hasEven := Some(numbers, func(n int) bool {
		return n%2 == 0
	})

	if !hasEven {
		t.Error("Expected to find even numbers")
	}
}

func TestEvery(t *testing.T) {
	numbers := []int{2, 4, 6, 8}
	allEven := Every(numbers, func(n int) bool {
		return n%2 == 0
	})

	if !allEven {
		t.Error("Expected all numbers to be even")
	}

	mixed := []int{1, 2, 3, 4}
	allEvenMixed := Every(mixed, func(n int) bool {
		return n%2 == 0
	})

	if allEvenMixed {
		t.Error("Expected not all numbers to be even")
	}
}

func TestPushPop(t *testing.T) {
	slice := []int{1, 2, 3}

	newLen := Push(&slice, 4, 5)
	if newLen != 5 {
		t.Errorf("Expected length 5, got %d", newLen)
	}

	element, ok := Pop(&slice)
	if !ok || element != 5 {
		t.Errorf("Expected to pop 5, got %d", element)
	}

	if len(slice) != 4 {
		t.Errorf("Expected length 4 after pop, got %d", len(slice))
	}
}

func TestIncludes(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}

	if !Includes(numbers, 3) {
		t.Error("Expected to include 3")
	}

	if Includes(numbers, 10) {
		t.Error("Expected not to include 10")
	}
}

func TestIndexOf(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}

	idx := IndexOf(numbers, 3)
	if idx != 2 {
		t.Errorf("Expected index 2, got %d", idx)
	}

	idx = IndexOf(numbers, 10)
	if idx != -1 {
		t.Errorf("Expected index -1, got %d", idx)
	}
}

func TestReverse(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	Reverse(numbers)

	expected := []int{5, 4, 3, 2, 1}
	for i, v := range numbers {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestConcat(t *testing.T) {
	slice1 := []int{1, 2}
	slice2 := []int{3, 4}
	slice3 := []int{5, 6}

	result := Concat(slice1, slice2, slice3)

	if len(result) != 6 {
		t.Errorf("Expected length 6, got %d", len(result))
	}

	expected := []int{1, 2, 3, 4, 5, 6}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestJoin(t *testing.T) {
	words := []string{"Hello", "World", "from", "Go"}
	result := Join(words, " ")

	expected := "Hello World from Go"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
