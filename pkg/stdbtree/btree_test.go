package stdbtree

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

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
	err := b.checkInvariances()
	require.NoError(t, err, testInfo)
	require.Equal(t, 0, b.len, testInfo)

	for _, num := range nums {
		prev := b.insert(num)
		require.Nil(t, prev, testInfo)
	}
	err = b.checkInvariances()
	require.NoError(t, err, testInfo)
	require.Equal(t, N, b.len, testInfo)

	// reinsert N items
	for _, num := range nums {
		prev := b.insert(num)
		require.NotNil(t, prev, testInfo)
		require.Equal(t, equal, prev.compare(num))
	}
	err = b.checkInvariances()
	require.NoError(t, err, testInfo)
	require.Equal(t, N, b.len, testInfo)

	// search for N items that we know are present
	for _, num := range nums {
		found := b.search(num)
		require.NotNil(t, found, testInfo)
		require.Equal(t, equal, num.compare(found))
	}
	err = b.checkInvariances()
	require.NoError(t, err, testInfo)
	require.Equal(t, N, b.len, testInfo)
}
