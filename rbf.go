package rbf


const NUM_TREES = 20
const TREE_SIZE = 1 << 25 // roughly 32M

const LEAF_SIZE = 64

const NUM_BITS = 32 // we want to store some shorts but also some ints, so need 4 bytes
const HIGH_BIT = 31 // low bit is 0, high bit is 31
const HIGH_BIT_1 = int32(-1) << HIGH_BIT

const NUM_FEATURES = int32(37 * 37 + 40)
const MAX_FEATURE = int32(NUM_FEATURES - 1)
const NUM_FEATURES_TO_COMPARE = int32(40)

const MAX_FEATURE_VALUE = 255 // openaddresses data has max followgram count ~200


const MAX_DEPTH = 25

type RandomBinaryTree struct {
    rowIndex   []int32
    treeFirst  []int32
    treeSecond []int32
    // TODO: REMOVE THESE
    numInternalNodes int32
    numLeaves        int32
}


type RandomBinaryForest struct {
    // TODO: make this private
    Trees               []RandomBinaryTree
}


var TEST_TRAINING_ADDRS []string

func init() {
    TEST_TRAINING_ADDRS = []string{"aaa", "abc"}
}

func MakeTestTree() RandomBinaryTree {
    // The test tree looks like it was "trained" on the strings "aaa" and "abc".
    // root node:
    //   treeFirst[0]: i.e. split on the 0 feature (i.e. "aa")
    //   treeSecond[0]: split-value 1
    // left child:
    //   treeFirst[1]: (leaf) 1 ("abc")  (actually HIGH_BIT_1 ^ 1)
    //   treeSecond[1]: (leaf) 2         (actualy HIGH_BIT_1 ^ 2)
    // right child:
    //   treeFirst[2]: (leaf) 0 ("aaa")  (actually HIGH_BIT_1 ^ 0)
    //   treeSecond[2]: (leaf) 1         (actually HIGH_BIT_1 ^ 1)
    rowIndex := []int32{0, 1}
    treeFirst := []int32{0, HIGH_BIT_1 ^ 1, HIGH_BIT_1 ^ 0}
    treeSecond := []int32{1, HIGH_BIT_1 ^ 2, HIGH_BIT_1 ^ 1}
    return RandomBinaryTree{rowIndex, treeFirst, treeSecond, 0, 0}
}
