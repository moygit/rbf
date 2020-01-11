package rbf

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestReadAndWriteForest(t *testing.T) {
	// given a test RBF
	rbf := RandomBinaryForest{[]RandomBinaryTree{NewTestTree(), NewTestTree()}}
	var builder strings.Builder
	// when we write it and read it back
	rbf.WriteToWriter(&builder)
	serialized := builder.String()
	rbfIn := ReadForestFromReader(strings.NewReader(serialized))
	// then the two forests should be identical
	if !reflect.DeepEqual(rbfIn, rbf) {
		t.Errorf("deserialized forest not the same as original forest")
	}
}
