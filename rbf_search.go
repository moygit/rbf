package rbf

import "github.com/Workiva/go-datastructures/bitarray"

// Query the forest for a single point. Returns indices into the training feature-array
// (since the caller/wrapper might have different things they want to do with this).
// Since multiple trees might return the same result points, this method dedups the results.
func (forest RandomBinaryForest) FindPointDedupResults(queryPoint []byte) []uint64 {
	// query each tree and get results (indices into training feature-array)
	resultIndices := bitarray.NewSparseBitArray()
	for _, tree := range forest.Trees {
		treeResultIndices := tree.findPoint(queryPoint)
		for _, index := range treeResultIndices {
			resultIndices.SetBit(uint64(index))
		}
	}
	return resultIndices.ToNums()
}

// Query the forest for a single point.
// Returns: for each tree in the forest, a slice of indices into the training feature-array
// (since the caller/wrapper might have different things they want to do with this).
// This method does not dedup results -- each tree's results are a different slice.
func (forest RandomBinaryForest) FindPointAllResults(queryPoint []byte) (count int, results [][]int32) {
	results = make([][]int32, len(forest.Trees))
	for i, tree := range forest.Trees {
		results[i] = tree.findPoint(queryPoint)
		count += len(results[i])
	}
	return
}

// A "point" is a feature-array. Search for one point in this tree.
func (tree RandomBinaryTree) findPoint(queryPoint []byte) []int32 {
	arrayPos := int32(0)
	first := tree.treeFirst[arrayPos]
	// the condition checks if it's an internal node (== 0) or a leaf (== -1):
	for first>>high_bit == 0 {
		// internal node, so first (the entry in tree.treeFirst) is a feature-number and
		// the entry in tree.treeSecond is the feature-value at which to split:
		if int32(queryPoint[first]) <= tree.treeSecond[arrayPos] {
			arrayPos = (2 * arrayPos) + 1
		} else {
			arrayPos = (2 * arrayPos) + 2
		}
		first = tree.treeFirst[arrayPos]
	}
	// found a leaf; get values and return
	indexStart, indexEnd := high_bit_1^first, high_bit_1^tree.treeSecond[arrayPos]
	return tree.rowIndex[indexStart:indexEnd]
}
