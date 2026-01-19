package gotool

import "time"

func BubbleSort(s []int) []int {
	for i := 0; i < len(s); i++ {
		time.Sleep(time.Millisecond)
		for j := 1; j < len(s); j++ {
			if s[j-1] > s[j] {
				s[j-1], s[j] = s[j], s[j-1]
			}
		}
	}
	return s
}
