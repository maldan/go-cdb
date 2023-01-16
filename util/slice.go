package util

func FilterSlice[T any](slice []T, filter func(*T) bool) []T {
	filtered := make([]T, 0)

	for i := 0; i < len(slice); i++ {
		if filter(&slice[i]) {
			filtered = append(filtered, slice[i])
		}
	}

	return filtered
}

func Unique[T comparable](slice []T) []T {
	keys := make(map[T]bool)
	list := make([]T, 0)

	for i := 0; i < len(slice); i++ {
		if _, value := keys[slice[i]]; !value {
			keys[slice[i]] = true
			list = append(list, slice[i])
		}
	}

	return list
}
