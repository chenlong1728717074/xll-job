package util

type String struct {
	value string
}

func (s String) Equal(other interface{}) bool {
	if s2, ok := other.(String); ok {
		return s.value == s2.value
	}
	return false
}
