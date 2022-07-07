package config

type confCollection[T comparable] []T

func newCollection[T comparable](v ...T) confCollection[T] {
	return v
}

func (m confCollection[T]) IndexOf(el T) int {
	for i, x := range m {
		if x == el {
			return i
		}
	}
	return -1
}
