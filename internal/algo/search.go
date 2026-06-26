package algo

func BinarySearch(arr []int, target int) (int, bool) {
	left := 0
	right := len(arr) - 1
	for left <= right {
		mid := left + (right - left) / 2
		if arr[mid] == target {
			return mid, true
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return -1, false
}