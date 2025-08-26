package cmp

type Equatable[T any] interface {
	Equal(T) bool
}

// SlicesEqualUnordered compares two slices of Equatable items for equality, regardless of order.
func SlicesEqualUnordered[T Equatable[T]](a, b []T, getKeyFunc func(T) string) bool {
	if len(a) != len(b) {
		return false
	}

	keyValue := make(map[string]T, len(a))
	keyCount := make(map[string]int, len(a))

	for _, v := range a {
		key := getKeyFunc(v)
		keyValue[key] = v
		keyCount[key]++
	}

	for _, v := range b {
		key := getKeyFunc(v)
		if keyCount[key] == 0 || !keyValue[key].Equal(v) {
			return false
		}
		keyCount[key]--
	}
	return true
}

func SliceEqualOrdered[T Equatable[T]](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if !v.Equal(b[i]) {
			return false
		}
	}

	return true
}

func PtrEqual[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
