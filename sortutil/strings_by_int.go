package sortutil

import "strconv"

type StringsByInt []string

func (b StringsByInt) Len() int {
	return len(b)
}

func (b StringsByInt) Less(i, j int) bool {
	iValue, _ := strconv.Atoi(b[i])
	jValue, _ := strconv.Atoi(b[j])
	return iValue < jValue
}

func (b StringsByInt) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
