package rbf


type RandomBinaryTree struct {
    // We have arrays of arrays of features. Instead of expensively moving those rows around when
    // sorting and partitioning we have an index into those and move the index elements around.
    // Lookups will be slightly more slower but we'll save time overall.
    rowIndex   []int32

    // Each tree node is a pair. For speed and space efficiency we'll store the tree in 2 arrays
    // using the standard trick for storing a binary tree in an array (with indexing starting at 0,
    // left child of n goes in 2n+1, right child goes in 2n+2). The pairs are either:
    // - if it's an internal node: the feature number and the value at which to split the feature
    // - if it's a leaf node: start and end indices in the rowIndex array; that view in the rowIndex
    //   array tells us the indices of rows in the original training set that are in this leaf
    // We distinguish the two cases by doing some bit-arithmetic.
    // 1. Yes, I know this is ugly, but the alternative is to have a whole 'nother pair of large arrays.
    // 2. Yes, I considered using hashmaps instead, but they're much slower (expected) and also take
    //    WAY more memory (which surprised me).
    treeFirst  []int32
    treeSecond []int32

    // TODO: THESE ARE FOR DEBUGGING AND WILL EVENTUALLY GO AWAY
    numInternalNodes int32
    numLeaves        int32
}


type RandomBinaryForest struct {
    Trees []RandomBinaryTree
}


// See comments above (in RandomBinaryTree definition) on ugly bit arithmetic for speed
const high_bit = 31 // low bit is 0, high bit is 31 because we're using int32s
const high_bit_1 = int32(-1) << high_bit
const max_feature_value = 255 // openaddresses data has max followgram count ~200



// It helps to have test values accessible from wrappers and sub-packages,
// so these need to be defined here.
var TEST_TRAINING_ADDRS []string

func init() {
    TEST_TRAINING_ADDRS = []string{"aaa", "abc"}
}

// Create a dummy tree for testing.
//
// The test tree looks like it was "trained" on the strings "aaa" and "abc".
//   root node:
//     treeFirst[0]: i.e. split on the 0 feature (i.e. "aa")
//     treeSecond[0]: split-value 1
//   left child:
//     treeFirst[1]: (leaf) 1 ("abc")  (actually high_bit_1 ^ 1)
//     treeSecond[1]: (leaf) 2         (actualy high_bit_1 ^ 2)
//   right child:
//     treeFirst[2]: (leaf) 0 ("aaa")  (actually high_bit_1 ^ 0)
//     treeSecond[2]: (leaf) 1         (actually high_bit_1 ^ 1)
func NewTestTree() RandomBinaryTree {
    rowIndex := []int32{0, 1}
    treeFirst := []int32{0, high_bit_1 ^ 1, high_bit_1 ^ 0}
    treeSecond := []int32{1, high_bit_1 ^ 2, high_bit_1 ^ 1}
    return RandomBinaryTree{rowIndex, treeFirst, treeSecond, 0, 0}
}


func check(e error) {
    if e != nil {
        panic(e)
    }
}
