package util

func ReverseSlice[T any](slice []T) []T {
	// 使用双指针交换法来倒序切片
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}
