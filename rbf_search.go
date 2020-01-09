package rbf


// A "point" is a feature-array. Search for one point in this tree.
func (tree RandomBinaryTree) FindPoint(queryPoint []byte) []int32 {
    arrayPos := int32(0)
    first := tree.treeFirst[arrayPos]
    // the condition checks if it's an internal node (== 0) or a leaf (== -1):
    for first >> high_bit == 0 {
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
    indexStart, indexEnd := high_bit_1 ^ first, high_bit_1 ^ tree.treeSecond[arrayPos]
    return tree.rowIndex[indexStart:indexEnd]
}


// A "point" is a feature-array. Search for one point in this forest.
// Return indices into the training feature-array (since the caller/wrapper might have
// different things they want to do with this).
func (forest RandomBinaryForest) FindPoint(queryPoint []byte) map[int32]bool {
    // query each tree and get results (indices into training feature-array)
    resultIndices := make(map[int32]bool)
    for _, tree := range forest.trees {
        treeResultIndices := tree.FindPoint(queryPoint)
        for _, index := range treeResultIndices {
            resultIndices[index] = true
        }
    }
    return resultIndices
}
