// Code generated by Constructor annotation processor. DO NOT EDIT.
// versions:
//		go: go1.18.3
//		go-annotation: 0.0.19
//		Constructor: 1.0.0

package utils

func NewNode[K comparable, V any]() *Node[K, V] {
	returnValue := &Node[K, V]{
		key:   []K{},
		nodes: make(map[K]*Node[K, V]),
		value: []V{},
	}

	return returnValue
}
