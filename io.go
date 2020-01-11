package rbf


import (
    "encoding/binary"
    "io"
    "reflect"
    "unsafe"
)


//######################################################################################################################
// Read/write int slices UNSAFELY.
//######################################################################################################################
// The godoc for the `binary` package says "This package favors simplicity over efficiency."
// They weren't kidding. Reading an int slice with binary.Read(...) takes 20x as long as reading
// a byte slice with bufio.Writer.Read(...) and then unsafely converting it to an int slice.
// Writing unsafely is "only" an order of magnitude faster.
// Conversion code h/t: https://stackoverflow.com/questions/11924196/convert-between-slices-of-different-types
const sizeof_int32 = 4 // bytes
func unsafelyReadIntSlice(reader io.Reader, int32sToRead int32) []int32 {
    // Read bytes from reader
    bytesToRead := 4 * int32sToRead
    byteSlice := make([]byte, bytesToRead)
    if n, err := reader.Read(byteSlice); n != int(bytesToRead) || err != nil {
        panic(err)
    }

    // Get the slice header and change length and capacity of slice
    header := *(*reflect.SliceHeader)(unsafe.Pointer(&byteSlice))
    header.Len /= sizeof_int32
    header.Cap /= sizeof_int32

    // Convert slice header to a []int32
    return *(*[]int32)(unsafe.Pointer(&header))
}
func unsafelyWriteIntSlice(writer io.Writer, intSlice []int32) {
    bytesToWrite := 4 * len(intSlice)

    // Get the slice header and change length and capacity of slice
    header := *(*reflect.SliceHeader)(unsafe.Pointer(&intSlice))
    header.Len *= sizeof_int32
    header.Cap *= sizeof_int32

    // Convert slice header to a []byte and write
    byteSlice := *(*[]byte)(unsafe.Pointer(&header))
    if n, err := writer.Write(byteSlice); n != bytesToWrite || err != nil {
        panic(err)
    }
}
//######################################################################################################################


func readTreeFromReader(reader io.Reader) RandomBinaryTree {
    // read lengths first so we can build slices to read
    var lenRowIndex, lenTreeFirst int32
    binary.Read(reader, binary.LittleEndian, &lenRowIndex)
    binary.Read(reader, binary.LittleEndian, &lenTreeFirst)

    // now read slices
    rowIndex := unsafelyReadIntSlice(reader, lenRowIndex)
    treeFirst := unsafelyReadIntSlice(reader, lenTreeFirst)
    treeSecond := unsafelyReadIntSlice(reader, lenTreeFirst)

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
    unsafelyWriteIntSlice(writer, tree.rowIndex)
    unsafelyWriteIntSlice(writer, tree.treeFirst)
    unsafelyWriteIntSlice(writer, tree.treeSecond)
    binary.Write(writer, binary.LittleEndian, int32(tree.numInternalNodes))
    binary.Write(writer, binary.LittleEndian, int32(tree.numLeaves))
}


func ReadForestFromReader(reader io.Reader) RandomBinaryForest {
    // first read the number of trees
    var numTrees int32
    err := binary.Read(reader, binary.LittleEndian, &numTrees)
    check(err)

    // and now read that each tree
    trees := make([]RandomBinaryTree, numTrees)
    for i := int32(0); i < numTrees; i++ {
        trees[i] = readTreeFromReader(reader)
    }

    return RandomBinaryForest{trees}
}


func (forest RandomBinaryForest) WriteToWriter(writer io.Writer) {
    // first write the number of trees
    err := binary.Write(writer, binary.LittleEndian, int32(len(forest.Trees)))
    check(err)
    // and now write each tree
    for _, tree := range forest.Trees {
        tree.writeToWriter(writer)
    }
}
