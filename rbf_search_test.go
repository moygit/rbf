package rbf


import (
    "testing"
)


func TestFindPoint(t *testing.T) {
    // given:
    trees := []RandomBinaryTree{NewTestTree(), NewTestTree()}
    forest := RandomBinaryForest{trees}
    queryPoint := []byte{6, 0, 0, 0, 0, 0}  // initial slice of followgrams for "aaaa"
    // when:
    queryResultIndices := forest.FindPoint(queryPoint)
    // then:
    if len(queryResultIndices) != 1 || !queryResultIndices[0] {
        t.Errorf("queryResultIndices == %v, expected %v", queryResultIndices, map[int32]bool{0: true})
    }
}
