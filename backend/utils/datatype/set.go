package datatype

type Set[T comparable] map[T] struct {}

func NewSet[T comparable]() Set[T] {
    return make(Set[T])
}

func (s Set[T]) Add(value T) {
    s[value] = struct{}{}
}

func (s Set[T]) Remove(value T) {
    delete(s, value)
}

func (s Set[T]) Size() int {
    return len(s)
}

func (s Set[T]) Contains(value T) bool {
	_, found := s[value]
	return found
}

func (s Set[T]) Values() []T {
    values := []T{}
    for key := range s {
        values = append(values, key)
    }
    return values
}