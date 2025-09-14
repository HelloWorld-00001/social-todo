package helper

func ListIntToInt32(list []int) []int32 {
	result := make([]int32, len(list))
	for i, v := range list {
		result[i] = int32(v)
	}
	return result
}

func ListInt32ToInt(list []int32) []int {
	result := make([]int, len(list))
	for i, v := range list {
		result[i] = int(v)
	}
	return result
}

func MapInt32ToInt(mp map[int32]int32) map[int]int {
	res := make(map[int]int, len(mp))
	for k, v := range mp {
		res[int(k)] = int(v)
	}
	return res
}

func MapIntToInt32(mp map[int]int) map[int32]int32 {
	res := make(map[int32]int32, len(mp))
	for k, v := range mp {
		res[int32(k)] = int32(v)
	}
	return res
}
