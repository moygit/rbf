package rbf


import (
    // "bytes"
    // "fmt"
    "os"
    // "reflect"
    "testing"
)


var TRAINING_ADDRS []string


func TestMain(m *testing.M) {
    TRAINING_ADDRS = []string{"aaa", "abc"}
    os.Exit(m.Run())
}


func makeTestTree() RandomBinaryTree {
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

func makeTestForest() RandomBinaryForest {
    addrs := []string{"aaa", "abc"}
    trees := []RandomBinaryTree{makeTestTree(), makeTestTree()}
    return RandomBinaryForest{addrs, trees}
}


func TestFindStringWithSimilarities(t *testing.T) {
    // given:
    queryString := "aaaa"
    forest := makeTestForest()
    // when:
    queryResults := forest.FindStringWithSimilarities(queryString)
    // then:
    if queryResults[0].Result != "aaa" {
        t.Errorf("queryResults[0].Result == %s, expected %s", queryResults[0].Result, TRAINING_ADDRS[0])
    }
    if queryResults[1].Result != "aaa" {
        t.Errorf("queryResults[1].Result == %s, expected %s", queryResults[1].Result, TRAINING_ADDRS[0])
    }
}

