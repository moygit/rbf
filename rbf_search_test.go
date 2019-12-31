package rbf


import (
    "testing"

    "rbf/features"
)


func makeTestForest() RandomBinaryForest {
    addrs := []string{"aaa", "abc"}
    trees := []RandomBinaryTree{MakeTestTree(), MakeTestTree()}
    featureSetConfigs := []features.FeatureSetConfig{ features.DefaultFollowgrams }
    return NewRBF(addrs, trees, featureSetConfigs)
}


func TestFindPoint(t *testing.T) {
    // given:
    forest := makeTestForest()
    queryPoint := []byte{6, 0, 0, 0, 0, 0}  // initial slice of features for "aaaa"
    // when:
    queryResultIndices := forest.FindPoint(queryPoint)
    // then:
    if len(queryResultIndices) != 1 || !queryResultIndices[0] {
        t.Errorf("queryResultIndices == %v, expected %v", queryResultIndices, map[int32]bool{0: true})
    }
}
