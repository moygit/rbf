package rbf


const NUM_TREES = 20
const TREE_SIZE = 1 << 25 // roughly 32M

const LEAF_SIZE = 8

const NUM_BITS = 32 // we want to store some shorts but also some ints, so need 4 bytes
const HIGH_BIT = 31 // low bit is 0, high bit is 31
const HIGH_BIT_1 = int32(-1) << HIGH_BIT

const NUM_FEATURES = 37 * 37
const MAX_FEATURE = NUM_FEATURES - 1
const SQRT_NUM_FEATURES = 37

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
    trainingStrings     []string
    // TODO: make this private
    Trees               []RandomBinaryTree
}
