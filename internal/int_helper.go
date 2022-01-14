package internal

import "strconv"

func StringArray2IntArray(in []string) ([]int, error) {
	var result []int
	var err error

	for _, i := range in {
		j, err := strconv.Atoi(i)
		if err != nil {
			break
		}
		result = append(result, j)
	}

	return result, err
}
