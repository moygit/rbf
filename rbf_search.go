package rbf

import (
    // "encoding/binary"
    // "fmt"
    // "io/ioutil"
    // "os"
    // "math/rand"
    "sort"

    "github.com/adrg/strutil"
    "github.com/adrg/strutil/metrics"

    // for logging only:
    // "log"
    // "os"
)


const NUM_RESULTS_TO_RETURN = 1
const CONSIDER_SIMILARITY_THRESHOLD = 0.4
var SIMILARITY_FUNC func(string, string) float64


type ResultSimilarityPair struct {
    // TODO: make these private
    Result string
    Similarity float64
}


func init() {
    // SIMILARITY_FUNC = rbfstrings.Jaccard2gramSimilarity
    lev := metrics.NewLevenshtein()
    SIMILARITY_FUNC = func(s1, s2 string) float64 {
            return strutil.Similarity(s1, s2, lev)
        }
}


// A "point" is a feature-array. Search for one point in this tree.
func (tree RandomBinaryTree) FindPoint(queryPoint []byte) []int32 {
    arrayPos := int32(0)
    first := tree.treeFirst[arrayPos]
    // the condition checks if it's an internal node (== 0) or a leaf (== 1):
    for !(first >> HIGH_BIT == -1) {
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
    indexStart, indexEnd := HIGH_BIT_1 ^ first, HIGH_BIT_1 ^ tree.treeSecond[arrayPos]
    return tree.rowIndex[indexStart:indexEnd]
}


// Search for the given input string in the given forest. Then get
// matching strings from the given list of training-strings, and get
// each result string's distance from the input string.
// Return the (upto) NUM_RESULTS_TO_RETURN best results.
func (forest RandomBinaryForest) FindStringWithSimilarities(queryString string) (ResultSimilarityPair, int, int) {
//func (forest RandomBinaryForest) FindStringWithSimilarities(queryString string) []ResultSimilarityPair {
    queryPoint := forest.calculateFeatures(queryString)

totalCount := 0
    // query each tree and get results (indices into forest.trainingStrings)
    resultIndices := make(map[int32]bool)
    for _, tree := range forest.Trees {
        treeResultIndices := tree.FindPoint(queryPoint)
        for _, index := range treeResultIndices {
            resultIndices[index] = true
totalCount += 1
        }
    }
distinctCount := len(resultIndices)

    // get strings and similarities
    results := make([]ResultSimilarityPair, 0)
    for index := range resultIndices {
        resultString := forest.trainingStrings[index]
        resultSimilarity := SIMILARITY_FUNC(queryString, resultString)
        // resultSimilarity := 0.0
        // if resultSimilarity > CONSIDER_SIMILARITY_THRESHOLD {
            results = append(results, ResultSimilarityPair{resultString, resultSimilarity})
        // }
    }

    // get the most similar results and return as many as we can, up to NUM_RESULTS_TO_RETURN
    sort.Slice(results, func (i int, j int) bool { return results[i].Similarity > results[j].Similarity })
    count := len(results)
    if count > NUM_RESULTS_TO_RETURN {
        count = NUM_RESULTS_TO_RETURN
    }
    // return results[:count]
return results[0], totalCount, distinctCount
}
