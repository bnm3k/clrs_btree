package main

import "fmt"

const greaterThan = 1
const equal = 0
const lessThan = -1

type item interface {
	compare(item) int
}

type node struct {
	isLeaf   bool
	items    []item
	children []*node
}

func newNode(t int, isLeaf bool) *node {
	items := make([]item, 0, 2*t-1)
	var children []*node = nil
	if !isLeaf { // if is internal
		children = make([]*node, 0, 2*t)
	}
	return &node{
		isLeaf:   isLeaf,
		items:    items,
		children: children,
	}
}

func (n *node) search(item item) item {
	for i := 0; i < len(n.items); i++ {
		switch item.compare(n.items[i]) {
		case greaterThan:
			continue
		case equal:
			return n.items[i]
		case lessThan:
			return n.children[i].search(item)
		}
	}
	if n.isLeaf {
		return nil
	}
	lastChildIndex := len(n.children) - 1
	return n.children[lastChildIndex].search(item)
}

type btree struct {
	t    int
	root *node
}

// t is the minimum degree a node is allowed to have.
// Every node must have t <= children <= 2t
// Exceptions: the root node may have less than t children.
// Every node must have t-1 <= keys  <= 2t - 1.
// Exceptions: the root node may have less than t-1 keys.
// t must be >= 2.
func newBTree(t int) *btree {
	if t < 2 {
		panic("invalid minimum degree for btree, t must be >= 2")
	}
	x := newNode(t, true)
	return &btree{
		t:    t,
		root: x,
	}
}

func main() {
	fmt.Println("hello world")
}
