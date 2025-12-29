package collection

/*
Returns a new slice with the transformed items.

If the slice is nil, returns nil.
*/
func Map[T any, U any](items []T, mapper func(T) U) []U {
	if items == nil {
		return nil
	}

	result := make([]U, len(items))
	for i, item := range items {
		result[i] = mapper(item)
	}

	return result
}

/*
Returns a new slice with the filtered items.

If the slice is nil, returns nil.
*/
func Filter[T any](items []T, predicate func(T) bool) []T {
	if items == nil {
		return nil
	}

	result := make([]T, 0)
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

/*
Returns true if every item in the slice satisfies the predicate.

If the slice is nil, returns true.
*/
func Every[T any](items []T, predicate func(T) bool) bool {
	if items == nil {
		return true
	}

	for _, item := range items {
		if !predicate(item) {
			return false
		}
	}

	return true
}
