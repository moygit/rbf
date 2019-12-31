package rbf


// A "random binary forest" is a hybrid between kd-trees and random forests.
// 
// We build an ensemble of roughly-binary search trees, with each tree being built as
// follows: pick a random subset of features at each split, look for the "best"
// feature, split on that feature, and then recurse.
// 
// For nearest-neighbor problems I've found that this performs significantly better
// than kd-trees (roughly as much better than plain kd-trees as random forests are
// better than plain decision-trees).
// 
// It's theoretically similar to a minhash except that the set of hashes is different
// for different subsets of input strings. If a node has two children then the node
// represents one hash, and the two children represent different hashes for different
// subsets of the input. The advantage over minhashes is that it's easier to have
// different *types* of features with an RBF than with a minhash.


import (
    "fmt"
    "os"
    "math"
    "math/rand"
    "sort"
    "sync"
    // "time"

    features "rbf/features"

    // for logging only:
    "log"
    // "os"
)


const LOG_FILENAME = "train.log"

var treeStatsFile *os.File
var logger *log.Logger

func init() {
    file, _ := os.OpenFile(LOG_FILENAME, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
    logger = log.New(file, "train_rbf ", log.Ldate|log.Ltime|log.Lshortfile)
treeStatsFile, _ = os.Create("tree_stats.txt")
}


func check(e error) {
    if e != nil {
        panic(e)
    }
}


// Convert a feature column into a feature-frequency histogram.
// Since our features are in the range [0, 255] (actually roughly [0, 200]), statistics will
// be easier if instead of passing around an array like [0, 1, 1, 1, 5, 5, 0, 0, 0, ...] (with
// millions of values) we instead pass around a frequency array like (for the above example)
// [4 (i.e. 4 0's), 3 (i.e. 3 1's), 0, 0, 0, 2].
// The followgram entries for the reference openaddresses dataset are all below 255, so we use 8 bits.
// Returns: for feature `feature_num`:
// - the frequency of each integer value in [0, 255]
// - the sum of all feature values (i.e. the weighted sum over the frequency array)
func getSingleFeatureFrequencies(rowIndex []int32, featureArray [][]byte, featureNum, indexStart, indexEnd int32) ([]int32, int32) {
    counts := make([]int32, MAX_FEATURE_VALUE+1)
    var weightedTotal int32 = 0
    for rowNum := indexStart; rowNum < indexEnd; rowNum++ {
        featureValue := featureArray[rowIndex[rowNum]][featureNum]
        counts[featureValue] += 1
        weightedTotal += int32(featureValue)
    }
    return counts, weightedTotal
    // Is it faster to calculate weightedTotal with additions inside the loop here,
    // or with 256 multiplications and additions on the counts list later?
    // Notes:
    // - "n-4" below is because when we get down close to the leaves we don't do this any more.
    // - This is all assuming it's a binary tree, which is obviously very approximate.
    // Calculations:
    // - additions inside loop:
    //   \\sum_{k=0}^{n-4} numNodes x numAdditions = \\sum_{k=0}^{n-4} 2^k 2^{n-k} = (n-3) * 2^n = n 2^n - 3 * 2^n
    // - 2 x 256 = 2^9 multiplications and additions on the counts list later:
    //   \\sum_{k=0}^{n-4} numNodes x 2 x 256 = \\sum_{k=0}^{n-4} 2^k 2^9 = 2^9 * (2^(n-3) - 1) = 2^6 2^n - 2^9
    // n is around 27 for USA (150 million), around 25 for France (25 million),
    // so for the full tree it's almost always faster to do them inside the above loop.
}


// Select a random subset of features and get the frequencies for those features.
func selectRandomFeaturesAndGetFrequencies(featureArray [][]byte, rowIndex []int32, indexStart, indexEnd int32) ([]int32, [][]int32, []int32) {
    featureSubset := make([]int32, NUM_FEATURES_TO_COMPARE)
    featureFrequencies := make([][]int32, NUM_FEATURES_TO_COMPARE)
    featureWeightedTotals := make([]int32, NUM_FEATURES_TO_COMPARE)

    var featureNum int32
    featuresAlreadySelected := make([]bool, NUM_FEATURES)
    for i := 0; i < NUM_FEATURES_TO_COMPARE; i++ {
        // get one that isn't already selected:
        for featureNum = int32(rand.Intn(NUM_FEATURES)); featuresAlreadySelected[featureNum]; featureNum = int32(rand.Intn(NUM_FEATURES)) {
        }
        featuresAlreadySelected[featureNum] = true
        featureSubset[i] = featureNum
        featureFrequencies[i], featureWeightedTotals[i] = getSingleFeatureFrequencies(rowIndex, featureArray, featureNum, indexStart, indexEnd)
    }
    return featureSubset, featureFrequencies, featureWeightedTotals
}


// Split a set of rows on one feature, trying to get close to the median but also maximizing variance.
//
// We want something as close to the median as possible so as to make the tree more balanced.
// And we want to calculate the "variance" about this split to compare features.
//
// CLEVERNESS ALERT (violating the "don't be clever" rule for speed):
// Except we'll actually use the mean absolute deviation instead of the variance as it's easier and
// better, esp since we're thinking of this in terms of Manhattan distance anyway. In fact, for our
// purposes it suffices to calculate the *total* absolute deviation, i.e. the total moment: we don't
// really need the mean since the denominator, the number of rows, is the same for all features that
// we're going to compare.
//
// The total moment to the right of some b, say for example b = 7.5, is
//     \\sum_{i=8}^{255} (i-7.5) * x_i = [ \\sum_{i=0}^255 (i-7.5) * x_i ] - [ \\sum_{i=0}^7 (i-7.5) x_i ]
// That second term is actually just -(the moment to the left of b), so the total moment
// (i.e. left + right) simplifies down to
//     \\sum_{i=0}^255 i x_i - \\sum_{i=0}^255 7.5 x_i + 2 \\sum_{i=0}^7 7.5 x_i - 2 \\sum_{i=0}^7 i x_i
// So we only need to track the running left-count and the running left-moment (w.r.t. 0), and then
// we can calculate the total moment w.r.t. median when we're done.
//
// Summary: Starting at 0.5 (no use starting at 0), iterate (a) adding to simple count, and (b)
// adding to left-side total moment. Stop as soon as the count is greater than half the total number
// of rows, and at that point we have a single expression for the total moment.
func splitOneFeature(featureHistogram []int32, totalZeroMoment int32, count int32) (float32, int32, int32) {
    fiftyPercentile := count / 2
    leftCount := featureHistogram[0]
    var pos, leftZeroMoment, thisItemCount, thisItemMoment int32
    for leftCount <= fiftyPercentile {
        pos += 1
        thisItemCount = featureHistogram[pos]
        thisItemMoment = thisItemCount * pos
        leftCount += thisItemCount
        leftZeroMoment += thisItemMoment
    }
    realPos := float32(pos) + 0.5   // want moment about e.g. 7.5, not 7 (using numbers in comment above)
    // See moment computation in comment above
    totalMoment := float32(totalZeroMoment) - (realPos * float32(count)) + (2 * ((realPos * float32(leftCount)) - float32(leftZeroMoment)))
    return totalMoment, pos, leftCount
}


type FeatureSplit struct {
    totalMoment float32
    splitValue  int32
    leftCount   int32
    featureNum  int
}

// Find the best of the given features, i.e. the one that has a split close to the median and has the highest variance.
// We only consider features that have a split between the 20th and 80th percentiles.
//
// Params:
// - featureFrequencies is an array giving the frequency (for that feature) of each integer value in [0, 255].
// - featureWeightedTotals is an array giving the weighted sum over the first array (same as the sum of all feature values).
// - totalCount is the number of rows which this iteration is looking at
// Returns: index and split-value of "best" feature, where
// - (index of) "best" = feature with greatest total absolute deviation about the median
// - split-value = (min value >= median)
const MIN_SPLIT_RATIO = 0.2
const MAX_SPLIT_RATIO = 0.8
func getBestFeature(featureFrequencies [][]int32, featureWeightedTotals []int32, totalCount int32) (int32, byte) {
    goodFeatureSplits := make([]FeatureSplit, len(featureFrequencies))
    goodCount := 0
    badFeatureSplits := make([]FeatureSplit, len(featureFrequencies))
    badCount := 0
    var featureSplits []FeatureSplit
    var count int

    for i, freq := range featureFrequencies {
        totalMoment, splitValue, leftCount := splitOneFeature(freq, featureWeightedTotals[i], totalCount)
        splitFrac := float64(leftCount) / float64(totalCount)
        if splitFrac > MIN_SPLIT_RATIO && splitFrac < MAX_SPLIT_RATIO {
            goodFeatureSplits[goodCount] = FeatureSplit{totalMoment, splitValue, leftCount, i}
            goodCount += 1
        } else {
            badFeatureSplits[badCount] = FeatureSplit{totalMoment, splitValue, leftCount, i}
            badCount += 1
        }
    }

    if goodCount > 0 {
        featureSplits = goodFeatureSplits
        count = goodCount
    } else {
        featureSplits = badFeatureSplits
        count = badCount
    }
    sort.Slice(featureSplits[:count],
               func(pos1, pos2 int) bool {
                   return featureSplits[pos1].totalMoment > featureSplits[pos2].totalMoment
               })
    bestFeature := featureSplits[0]
    return int32(bestFeature.featureNum), byte(bestFeature.splitValue)
}


type simpleFeatureSplit struct {
    splitDiff   float32
    splitValue  int32
    leftCount   int32
    featureNum  int
}

// From the given features find the one which splits closest to the median.
func getSimpleBestFeature(featureFrequencies [][]int32, featureWeightedTotals []int32, totalCount int32) (int32, byte) {
    featureSplits := make([]simpleFeatureSplit, len(featureFrequencies))

    for i, freq := range featureFrequencies {
        _, splitValue, leftCount := splitOneFeature(freq, featureWeightedTotals[i], totalCount)
        splitDiff := math.Abs(float64(leftCount - (totalCount - leftCount)))    // leftCount - rightCount
        featureSplits[i] = simpleFeatureSplit{float32(splitDiff), splitValue, leftCount, i}
    }

    sort.Slice(featureSplits,
               func(pos1, pos2 int) bool {
                   return featureSplits[pos1].splitDiff < featureSplits[pos2].splitDiff
               })
    bestFeature := featureSplits[0]
    return int32(bestFeature.featureNum), byte(bestFeature.splitValue)
}


// quicksort-type partitioning of rowIndex[indexStart..indexEnd) based on whether the
// feature `featureNum` is less-than-or-equal-to or greater-than splitValue
// pre-req: indexEnd - indexStart is at least 2
func quickPartition(rowIndex []int32, featureArray [][]byte, indexStart, indexEnd, featureNum int32, splitValue byte) int32 {
    for i, j := indexStart, indexEnd-1; i < j; {
        for i < indexEnd && featureArray[rowIndex[i]][featureNum] <= splitValue {
            i += 1
        }
        for j >= indexStart && featureArray[rowIndex[j]][featureNum] > splitValue {
            j -= 1
        }
        if i >= j {
            return i
        }
        rowIndex[i], rowIndex[j] = rowIndex[j], rowIndex[i]
    }
    return indexStart // should never get here unless passed illegal values (start >= end)
}


// Allocate space for the tree's component arrays and then
// call the recursive `calculateOneNode` function, which does the real training.
func trainOneTree(treeNum int, featureArray [][]byte) RandomBinaryTree {
    rowIndex := make([]int32, len(featureArray))
    for i := int32(0); i < int32(len(rowIndex)); i++ {
        rowIndex[i] = i
    }
    treeFirst := make([]int32, TREE_SIZE)
    treeSecond := make([]int32, TREE_SIZE)
    var numInternalNodes int32
    var numLeaves int32
    calculateOneNode(treeNum, featureArray, rowIndex, treeFirst, treeSecond, &numInternalNodes, &numLeaves, 0, int32(len(rowIndex)), 0, 0)
    return RandomBinaryTree{rowIndex, treeFirst, treeSecond, numInternalNodes, numLeaves}
}


func TrainForestWithFeatureArray(featureArray [][]byte) []RandomBinaryTree {
    // make and train trees:
    trees := make([]RandomBinaryTree, NUM_TREES)
    var wg sync.WaitGroup
    for i := 0; i < NUM_TREES; i++ {
        wg.Add(1)
        go func(j int) {
            defer wg.Done()
            trees[j] = trainOneTree(j, featureArray)
        }(i)
    }
    wg.Wait()
treeStatsFile.Close()
    return trees
}


//  // TODO: Fix this comment
//  // TODO: I DON'T KNOW HOW TO EXTEND RIGHT NOW because we're not using List or array.array or
//  //       numpy.array, we're using multiprocessing.sharedctypes.RawArray
//  func setArrays(arraypos, val1, val2):
//      # if pos >= len(treeFirst):  # the two arrays have the same dimensions
//      #     treeFirst.extend([0] * len(treeFirst))
//      #     treeSecond.extend([0] * len(treeSecond))
//      # if pos >= len(treeFirst):  # the two arrays have the same dimensions
//      #     # right child node is at 2n+1, so we can never need more than len(treeArray)+1 new nodes
//      #     np.append(treeFirst, [0] * (len(treeFirst) + 1))
//      #     np.append(treeSecond, [0] * (len(treeSecond) + 1))
//      treeFirst[pos] = val1
//      treeSecond[pos] = val2


// Calculate the split (or leaf) at one node (and its descendants).
// Params:
// - TODO: update this
// - indexStart and indexEnd: the view into rowIndex that we're considering right now
// - treeArrayPos: the position of this node in the tree arrays
// Guarantees:
// - Parallel calls to `calculateOneNode` will look at non-intersecting views.
// - Child calls will look at distinct sub-views of this view.
// - No two calls to `calculateOneNode` will have the same treeArrayPos
func calculateOneNode(treeNum int, featureArray [][]byte, rowIndex, treeFirst, treeSecond []int32,
                      numInternalNodes, numLeaves *int32,
                      indexStart, indexEnd int32, treeArrayPos int, depth int) {
    // logger.Printf("indexStart: %d, indexEnd: %d, treeArrayPos: %d\n", indexStart, indexEnd, treeArrayPos)
    if 2*treeArrayPos+2 >= len(treeFirst) {
    // Special termination condition to regulate depth.
        // logger.Printf("DEBUG:  depth check triggered, 2 * treeArrayPos + 1 >= len(treeFirst): 2 * %d + 1 >= %d\n", treeArrayPos, len(treeFirst))
        // logger.Printf("WARNING: reached max depth: 2 * treeArrayPos >= len(treeFirst) with indexStart and indexEnd: "+
            // "2 * %d >= %d with %d and %d\n", treeArrayPos, len(treeFirst), indexStart, indexEnd)
        treeFirst[treeArrayPos], treeSecond[treeArrayPos] = HIGH_BIT_1 ^ indexStart, HIGH_BIT_1 ^ indexEnd
        // TODO: REMOVE THIS!
        *numLeaves += 1
fmt.Fprintf(treeStatsFile, "%d,%d,%d,depth-based-leaf,%d,%d,%d,%d,%d,%d,\n", treeNum, treeArrayPos, depth, indexStart, indexEnd, indexEnd-indexStart, 0, 0, 0)
        return
    }

    if indexEnd-indexStart < LEAF_SIZE {
    // Not enough items left to split. Make a leaf.
        // logger.Printf("DEBUG: making leaf")
        treeFirst[treeArrayPos], treeSecond[treeArrayPos] = HIGH_BIT_1 ^ indexStart, HIGH_BIT_1 ^ indexEnd
        // TODO: REMOVE THIS!
        *numLeaves += 1
fmt.Fprintf(treeStatsFile, "%d,%d,%d,size-based-leaf,%d,%d,%d,%d,%d,%d,\n", treeNum, treeArrayPos, depth, indexStart, indexEnd, indexEnd-indexStart, 0, 0, 0)
    } else {
    // Not a leaf. Get a random subset of NUM_FEATURES_TO_COMPARE features, find the best one, and split this node.
    // TODO (not sure where): pick feature so that each side has at least a third of data, else don't bother splitting if below a threshold
    //      or look at more features or something
        // logger.Printf("DEBUG: splitting node")
        featureSubset, featureFrequencies, featureWeightedTotals := selectRandomFeaturesAndGetFrequencies(featureArray, rowIndex, indexStart, indexEnd)
        bestFeatureIndex, bestFeatureSplitValue := getSimpleBestFeature(featureFrequencies, featureWeightedTotals, indexEnd-indexStart)
        bestFeatureNum := featureSubset[bestFeatureIndex]
        indexSplit := quickPartition(rowIndex, featureArray, indexStart, indexEnd, bestFeatureNum, bestFeatureSplitValue)
        // logger.Printf("DEBUG: bestFeatureSplitValue: %d, bestFeatureNum: %d, indexSplit: %d\n", bestFeatureSplitValue, bestFeatureNum, indexSplit)

        if indexSplit == indexStart || indexSplit == indexEnd {
        // BUG TODO: the first set of features gave us a crappy split. What do we do?
            logger.Printf("DEBUG: bad split; feature-num: %d, count: %d", bestFeatureNum, indexEnd-indexStart)
        }
        treeFirst[treeArrayPos], treeSecond[treeArrayPos] = bestFeatureNum, int32(bestFeatureSplitValue)
fmt.Fprintf(treeStatsFile, "%d,%d,%d,internal,%d,%d,%d,%d,%d,%d,%s\n", treeNum, treeArrayPos, depth, indexStart, indexEnd, indexEnd-indexStart, indexSplit, bestFeatureNum, bestFeatureSplitValue, features.CHAR_REVERSE_MAP[bestFeatureNum])
        // TODO: REMOVE THIS!
        *numInternalNodes += 1
        calculateOneNode(treeNum, featureArray, rowIndex, treeFirst, treeSecond, numInternalNodes, numLeaves, indexStart, indexSplit, (2*treeArrayPos)+1, depth+1)
        calculateOneNode(treeNum, featureArray, rowIndex, treeFirst, treeSecond, numInternalNodes, numLeaves, indexSplit, indexEnd, (2*treeArrayPos)+2, depth+1)
    }
}
