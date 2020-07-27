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
	n        int // tracks no. of items in a node
	items    []item
	children []*node
}

// for debugging
func (n *node) toString() string {
	s := "{"
	s += fmt.Sprintf("isLeaf:%5v, ", n.isLeaf)
	s += fmt.Sprintf("n:%2d, ", n.n)
	s += fmt.Sprintf("items: %v", n.items[:n.n])
	if !n.isLeaf {
		s += "\n\n\t"
		s += n.children[0].toString()
		for i := 1; i <= n.n; i++ {
			s += "       "
			s += n.children[i].toString()
		}
		s += "\n\n"
	}
	s += "}"
	return s
}

func newNode(t int, isLeaf bool) *node {
	items := make([]item, 2*t-1)
	var children []*node = nil
	if !isLeaf { // if is internal
		children = make([]*node, 2*t)
	}
	return &node{
		isLeaf:   isLeaf,
		items:    items,
		children: children,
	}
}

func (n *node) search(item item) item {
	for i := 0; i < n.n; i++ {
		switch item.compare(n.items[i]) {
		case greaterThan:
			break
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

func (n *node) insertLeaf(newItem item) (prev item) {
	var i int
loop:
	for i = 0; i < n.n; i++ {
		curr := n.items[i]
		switch newItem.compare(curr) {
		case equal:
			prev = curr
			break loop
		case lessThan:
			copy(n.items[i+1:], n.items[i:])
			break loop
		}
	}
	n.items[i] = newItem
	if prev == nil { // i.e. is fresh insert
		n.n++
	}
	return
}

func (n *node) insert(t int, newItem item) (prev item) {
	if n.isLeaf {
		return n.insertLeaf(newItem)
	}
	var i int
loop:
	for i = 0; i < n.n; i++ {
		curr := n.items[i]
		switch newItem.compare(curr) {
		case equal:
			prev = curr
			n.items[i] = newItem
			return
		case lessThan:
			break loop
		}
	}
	c := n.children[i]
	if c.n == 2*t-1 {
		median := n.splitChild(t, i)
		switch newItem.compare(median) {
		case lessThan:
			// go to left child
		case equal:
			// replace
			prev = median
			n.items[i] = newItem
			return
		case greaterThan:
			// go to newly upped right child
			c = n.children[i+1]
		}
	}
	return c.insert(t, newItem)
}

func (n *node) splitChild(t int, i int) (median item) {
	// let y be the ith child of node n.
	y := n.children[i]
	median = y.items[t-1]

	// halve y and move the upper half to new node z
	z := newNode(t, y.isLeaf)
	copy(z.items, y.items[t:])
	z.n = t - 1
	y.n = t - 1
	if !y.isLeaf { // only internal nodes have children
		copy(z.children, y.children[t:])
	}

	// move median item up to parent (node n)
	copy(n.items[i+1:], n.items[i:])
	n.items[i] = median
	n.n++

	// add z as node n's child
	copy(n.children[i+2:], n.children[i+1:])
	n.children[i+1] = z
	return median
}

type btree struct {
	root *node
	t    int
	len  int
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

func (b *btree) search(item item) item {
	return b.root.search(item)
}

// for debugging
func (b *btree) checkInvariances() error {
	var traverseItems func(n *node, fn func(i item))
	traverseItems = func(n *node, fn func(i item)) {
		var i int
		for i = 0; i < n.n; i++ {
			// first traverse children if internal
			if !n.isLeaf {
				traverseItems(n.children[i], fn)
			}
			fn(n.items[i])
		}
		// traverse last child
		if !n.isLeaf {
			traverseItems(n.children[i], fn)
		}
	}

	// check that there are no duplicates and all items are in ascending order
	// this also implictly checks that for every key k, all the items at that subtree
	// are less than key k
	var items []item
	traverseItems(b.root, func(i item) {
		items = append(items, i)
	})
	for i := 1; i < len(items); i++ {
		switch items[i].compare(items[i-1]) {
		case equal:
			return fmt.Errorf("btree contains duplicate items: %v, %v", items[i-1], items[i])
		case lessThan:
			return fmt.Errorf("btree items not in sorted order (ascending)\n: %v comes before %v", items[i-1], items[i])
		}
	}

	// preOrder-ish traversal, ie traverse node then children
	var traverseNode func(n *node, fn func(n *node))
	traverseNode = func(n *node, fn func(n *node)) {
		fn(n)
		if !n.isLeaf {
			for i := 0; i < n.n+1; i++ {
				traverseNode(n.children[i], fn)
			}
		}
	}

	// check that all nodes have correct n
	if b.root.n > 2*b.t-1 {
		return fmt.Errorf("Root node has invalid n: %d", b.root.n)
	}
	var err error
	if !b.root.isLeaf {
		for i := 0; i < b.root.n+1; i++ {
			traverseNode(b.root.children[i], func(n *node) {
				if n.n < b.t-1 || n.n > 2*b.t-1 {
					err = fmt.Errorf("One of the nodes has invalid n: %d", n.n)
				}
			})
		}
	}
	if err != nil {
		return err
	}

	// check that all leaves are at same height
	var leafHeights []int
	var traverseHeight func(n *node, level int)
	traverseHeight = func(n *node, level int) {
		if n.isLeaf {
			leafHeights = append(leafHeights, level)
		} else {
			for i := 0; i <= n.n; i++ {
				traverseHeight(n.children[i], level+1)
			}
		}
	}
	traverseHeight(b.root, 1)
	height := leafHeights[0]
	for _, h := range leafHeights {
		if h != height {
			return fmt.Errorf("one of the leaf nodes does not have the same height as the rest: %d vs %d", h, height)
		}
	}
	return nil
}

func (b *btree) insert(item item) (prev item) {
	if b.root.n == (2*b.t - 1) {
		oldRoot := b.root
		b.root = newNode(b.t, false)
		b.root.children[0] = oldRoot
		b.root.splitChild(b.t, 0)
	}
	prev = b.root.insert(b.t, item)
	if prev == nil {
		b.len++
	}
	return
}

func insertAt(arr []int, index int, item int) []int {
	copy(arr[index+1:], arr[index:])
	arr[index] = item
	return arr
}

type numItem int

func (n numItem) compare(other item) int {
	otherNum, ok := other.(numItem)
	if !ok {
		panic("invalid item type for comparison")
	}
	if n < otherNum {
		return lessThan
	} else if n == otherNum {
		return equal
	}
	return greaterThan
}

func main() {
	b := newBTree(2)
	var n int = 300
	for i := 1; i <= 100; i++ {
		b.insert(numItem(n))
		n++
	}
	if err := b.checkInvariances(); err != nil {
		panic(err)
	}
}
