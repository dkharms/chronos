package gotool

func BubbleSort(s []int) []int {
	for i := 0; i < len(s); i++ {
		for j := 1; j < len(s); j++ {
			if s[j-1] > s[j] {
				s[j-1], s[j] = s[j], s[j-1]
			}
		}
	}
	return s
}
