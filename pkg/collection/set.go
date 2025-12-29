package collection

type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

func (s Set[T]) Add(item T) bool {
	_, ok := s[item]
	if ok {
		return false
	}

	s[item] = struct{}{}
	return true
}

func (s Set[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}

func (s Set[T]) Remove(item T) bool {
	_, ok := s[item]
	if !ok {
		return false
	}

	delete(s, item)
	return true
}

func (s Set[T]) Slice() []T {
	out := make([]T, 0, len(s))
	for k := range s {
		out = append(out, k)
	}
	return out
}
