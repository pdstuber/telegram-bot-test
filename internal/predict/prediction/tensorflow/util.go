package tensorflow

import "math"

func findMaxAndMaxIndex(probs []float32) (float32, int) {
	max := float32(math.MinInt32)
	maxIndex := -1

	for i := range probs {
		value := probs[i]
		if value > max {
			max = value
			maxIndex = i
		}
	}

	return max, maxIndex
}
func findClassWithMaxProbability(probs []float32, labels []Label) (string, float32) {
	max, maxIndex := findMaxAndMaxIndex(probs)

	return labels[maxIndex].ClassName, max
}
