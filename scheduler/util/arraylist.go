package util

type Entry interface {
	Equal(other interface{}) bool
}

type ArrayList[T Entry] struct {
	elements []T
	size     int64
	cap      int64
}

func NewArrayList[T Entry]() *ArrayList[T] {
	return &ArrayList[T]{
		elements: make([]T, 0),
	}
}

func (list *ArrayList[T]) Contains(item T) bool {
	for _, element := range list.elements {
		if element.Equal(item) {
			return true
		}
	}
	return false
}

func (list *ArrayList[T]) Add(item T) {
	list.elements = append(list.elements, item)
}

func (list *ArrayList[T]) Remove(item T) {
	list.elements = append(list.elements, item)
}

func (list *ArrayList[T]) Size() int64 {
	return list.size
}
