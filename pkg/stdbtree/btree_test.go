package stdbtree

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// for debugging/testing
func checkInvariances(b *btree) error {
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

func TestBtreeBasic(t *testing.T) {
	// check that t must b >= 2
	require.Panics(t, func() {
		newBTree(1)
	})

	// test parameters
	N := 300
	T := 2
	seedVal := time.Now().UnixNano()
	testInfo := fmt.Sprintf("[seedVal = %d, T = %d]", seedVal, T) // for replication

	// test values
	var nums []numItem
	for i := 0; i < N; i++ {
		nums = append(nums, numItem(i))
	}
	rand.Seed(seedVal)
	rand.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })

	// newBtree
	b := newBTree(T)
	require.NotNil(t, b)
	err := checkInvariances(b)
	require.NoError(t, err, testInfo)
	require.Equal(t, 0, b.len, testInfo)

	for _, num := range nums {
		prev := b.insert(num)
		require.Nil(t, prev, testInfo)
	}
	err = checkInvariances(b)
	require.NoError(t, err, testInfo)
	require.Equal(t, N, b.len, testInfo)

	// reinsert N items
	for _, num := range nums {
		prev := b.insert(num)
		require.NotNil(t, prev, testInfo)
		require.Equal(t, equal, prev.compare(num))
	}
	err = checkInvariances(b)
	require.NoError(t, err, testInfo)
	require.Equal(t, N, b.len, testInfo)

	// search for N items that we know are present
	for _, num := range nums {
		found := b.search(num)
		require.NotNil(t, found, testInfo)
		require.Equal(t, equal, num.compare(found))

	}
	err = checkInvariances(b)
	require.NoError(t, err, testInfo)
	require.Equal(t, N, b.len, testInfo)

}
