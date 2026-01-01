package collection

import "testing"

func TestMap(t *testing.T) {
	t.Run("nil slice returns nil", func(t *testing.T) {
		var items []int
		got := Map(items, func(v int) int { return v * 2 })
		if got != nil {
			t.Fatal("expected nil slice")
		}
	})

	t.Run("maps all items", func(t *testing.T) {
		items := []int{1, 2, 3}
		got := Map(items, func(v int) int { return v * 2 })

		want := []int{2, 4, 6}
		if len(got) != len(want) {
			t.Fatalf("expected len %d, got %d", len(want), len(got))
		}

		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("at index %d expected %d, got %d", i, want[i], got[i])
			}
		}
	})
}

func TestFilter(t *testing.T) {
	t.Run("nil slice returns nil", func(t *testing.T) {
		var items []int
		got := Filter(items, func(v int) bool { return v%2 == 0 })
		if got != nil {
			t.Fatal("expected nil slice")
		}
	})

	t.Run("filters items correctly", func(t *testing.T) {
		items := []int{1, 2, 3, 4}
		got := Filter(items, func(v int) bool { return v%2 == 0 })

		want := []int{2, 4}
		if len(got) != len(want) {
			t.Fatalf("expected len %d, got %d", len(want), len(got))
		}

		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("at index %d expected %d, got %d", i, want[i], got[i])
			}
		}
	})
}

func TestEvery(t *testing.T) {
	t.Run("nil slice returns true", func(t *testing.T) {
		var items []int
		if !Every(items, func(v int) bool { return v > 0 }) {
			t.Fatal("expected true for nil slice")
		}
	})

	t.Run("all items satisfy predicate", func(t *testing.T) {
		items := []int{2, 4, 6}
		if !Every(items, func(v int) bool { return v%2 == 0 }) {
			t.Fatal("expected true when all items satisfy predicate")
		}
	})

	t.Run("some item does not satisfy predicate", func(t *testing.T) {
		items := []int{2, 3, 6}
		if Every(items, func(v int) bool { return v%2 == 0 }) {
			t.Fatal("expected false when any item fails predicate")
		}
	})
}
