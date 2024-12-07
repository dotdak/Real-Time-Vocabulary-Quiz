package utils

type HashSet[K comparable] struct {
	Data map[K]struct{}
}

func NewSet[K comparable](keys ...K) *HashSet[K] {
	set := &HashSet[K]{
		Data: make(map[K]struct{}),
	}

	set.Adds(keys...)

	return set
}

func (s *HashSet[K]) Add(k K) {
	s.Data[k] = struct{}{}
}

func (s *HashSet[K]) Adds(keys ...K) {
	for _, k := range keys {
		s.Data[k] = struct{}{}
	}
}

func (s *HashSet[K]) Has(k K) bool {
	_, ok := s.Data[k]
	return ok
}

func (s *HashSet[K]) Pop(k K) {
	delete(s.Data, k)
}

func (s *HashSet[K]) Len() int {
	return len(s.Data)
}

func (s *HashSet[K]) Keys() []K {
	out := make([]K, 0, s.Len())
	for k := range s.Data {
		out = append(out, k)
	}

	return out
}
