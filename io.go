package rbf


import (
    "encoding/binary"
    "io"
)


func readTreeFromReader(reader io.Reader) RandomBinaryTree {
    // read lengths first so we can build slices to read next
    var lenRowIndex, lenTreeFirst int32
    binary.Read(reader, binary.LittleEndian, &lenRowIndex)
    binary.Read(reader, binary.LittleEndian, &lenTreeFirst)

    // now read arrays
    rowIndex := make([]int32, lenRowIndex)
    treeFirst := make([]int32, lenTreeFirst)
    treeSecond := make([]int32, lenTreeFirst)
    binary.Read(reader, binary.LittleEndian, &rowIndex)
    binary.Read(reader, binary.LittleEndian, &treeFirst)
    binary.Read(reader, binary.LittleEndian, &treeSecond)

    // TODO: remove this
    // and finally read node counts
    var numInternalNodes, numLeaves int32
    binary.Read(reader, binary.LittleEndian, &numInternalNodes)
    binary.Read(reader, binary.LittleEndian, &numLeaves)
    return RandomBinaryTree{rowIndex, treeFirst, treeSecond, numInternalNodes, numLeaves}
}


func (tree RandomBinaryTree) writeToWriter(writer io.Writer) {
    binary.Write(writer, binary.LittleEndian, int32(len(tree.rowIndex)))
    binary.Write(writer, binary.LittleEndian, int32(len(tree.treeFirst)))
    binary.Write(writer, binary.LittleEndian, tree.rowIndex)
    binary.Write(writer, binary.LittleEndian, tree.treeFirst)
    binary.Write(writer, binary.LittleEndian, tree.treeSecond)
    binary.Write(writer, binary.LittleEndian, tree.numInternalNodes)
    binary.Write(writer, binary.LittleEndian, tree.numLeaves)
}


// TODO: error-handling
func ReadForestFromReader(reader io.Reader) RandomBinaryForest {
    // read numTrees first, then read that many trees
    var numTrees int32
    binary.Read(reader, binary.LittleEndian, &numTrees)

    trees := make([]RandomBinaryTree, numTrees)
    for i := int32(0); i < numTrees; i++ {
        trees[i] = readTreeFromReader(reader)
    }

    return RandomBinaryForest{trees}
}


func (forest RandomBinaryForest) WriteToWriter(writer io.Writer) {
    binary.Write(writer, binary.LittleEndian, int32(len(forest.Trees)))
    for _, tree := range forest.Trees {
        tree.writeToWriter(writer)
    }
}
