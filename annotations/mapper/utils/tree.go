package utils

/*
The tree algorithm should splice single nodes to one
FROM:
				a
			   /| \
		     /  |   \
		  /     |    \
		/		|	   \
       b		|		|
	   |		|		|
       d		e		f
	  /\		|	   / \
	/	\		|     /    \
	k    p		l	  n		m
	|			|			|
	x			y			z
TO:
				a
			   /| \
		     /  |   \
		  /     |    \
		/	   ely	   \
       bd	   			f
	  /\			   / \
	/	\		      /    \
   kx    p			  n	   mz
*/

// Node base tree structure
// @Constructor(type="pointer")
type Node[K comparable, V any] struct {
	key   []K               //@Init
	value []V               //@Init
	nodes map[K]*Node[K, V] //@Init
}

func (n *Node[K, V]) Add(key []K, value V) {
	if len(key) == 0 {
		n.value = append(n.value, value)
		return
	}
	k1, k2 := key[0], key[1:]
	node := n.getOrInitNode(k1)
	node.Add(k2, value)
}

func (n *Node[K, V]) Optimize() {
	if len(n.nodes) == 0 {
		return
	}
	if len(n.value) == 0 && len(n.nodes) == 1 {
		for _, node := range n.nodes {
			n.value = node.value
			n.key = append(n.key, node.key...)
			n.nodes = node.nodes
			n.Optimize()
			return
		}
	}
	for _, node := range n.nodes {
		node.Optimize()
	}
}

func (n *Node[K, V]) Execute(pre, post func([]K, []V)) {
	pre(n.key, n.value)
	defer post(n.key, n.value)
	for _, node := range n.nodes {
		node.Execute(pre, post)
	}
}

func (n *Node[K, V]) getOrInitNode(key K) *Node[K, V] {
	node, ok := n.nodes[key]
	if !ok {
		node = NewNode[K, V]()
		node.key = append(node.key, key)
		n.nodes[key] = node
	}
	return node
}
