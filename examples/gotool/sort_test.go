package gotool

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"
)

func TestBubbleSort(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		want []int
	}{
		{
			name: "dead-simple",
			s:    []int{3, 2, 1},
			want: []int{1, 2, 3},
		},
		{
			name: "one-inversion",
			s:    []int{1, 2, 4, 3, 5},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "indempotent",
			s:    []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !slices.IsSorted(BubbleSort(tt.s)) {
				t.Fail()
			}
		})
	}
}

func BenchmarkBubbleSort(b *testing.B) {
	generate := func(to, length int) []int {
		s := make([]int, 0, length)
		for range length {
			s = append(s, rand.Intn(to))
		}
		return s
	}

	for _, l := range r(7, 14) {
		exp := 1 << l
		b.Run(fmt.Sprintf("length=%d", exp), func(b *testing.B) {
			for b.Loop() {
				b.StopTimer()
				s := generate(1024, exp)
				b.StartTimer()

				BubbleSort(s)
			}
		})
	}
}

func r(from, to int) (s []int) {
	for i := from; i <= to; i++ {
		s = append(s, i)
	}
	return
}
