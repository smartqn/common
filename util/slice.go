package util

func InStrSlice(target string, slice []string) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == target {
			return true
		}
	}
	return false
}

func InIntSlice(target int, slice []int) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == target {
			return true
		}
	}
	return false

}

func InInt64Slice(target int64, slice []int64) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == target {
			return true
		}
	}
	return false

}
