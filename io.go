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
// Conversion code h/t:
// https://stackoverflow.com/questions/11924196/convert-between-slices-of-different-types
const sizeof_int32 = 4 // bytes
func unsafelyReadIntSlice(reader io.Reader, int32sToRead int32) []int32 {
    bytesToRead := 4 * int32sToRead
    byteSlice := make([]byte, bytesToRead)
    if n, err := reader.Read(byteSlice); n != int(bytesToRead) || err != nil {
        // TODO: handle the error!
    }

    // Get the slice header
    header := *(*reflect.SliceHeader)(unsafe.Pointer(&byteSlice))

    // Change length and capacity of slice
    header.Len /= sizeof_int32
    header.Cap /= sizeof_int32

    // Convert slice header to an []int32
    return *(*[]int32)(unsafe.Pointer(&header))
}
func unsafelyWriteIntSlice(writer io.Writer, intSlice []int32) {
    bytesToWrite := 4 * len(intSlice)

    // Get the slice header
    header := *(*reflect.SliceHeader)(unsafe.Pointer(&intSlice))

    // Change length and capacity of slice
    header.Len *= sizeof_int32
    header.Cap *= sizeof_int32

    // Convert slice header to an []int32
    byteSlice := *(*[]byte)(unsafe.Pointer(&header))
    if n, err := writer.Write(byteSlice); n != int(bytesToWrite) || err != nil {
        // TODO: handle the error!
    }
}
//######################################################################################################################


func readTreeFromReader(reader io.Reader) RandomBinaryTree {
    // read lengths first so we can build slices to read next
    var lenRowIndex, lenTreeFirst int32
    binary.Read(reader, binary.LittleEndian, &lenRowIndex)
    binary.Read(reader, binary.LittleEndian, &lenTreeFirst)

    // now read arrays
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
    binary.Write(writer, binary.LittleEndian, tree.numInternalNodes)
    binary.Write(writer, binary.LittleEndian, tree.numLeaves)
}


// TODO: error-handling
func ReadForestFromReader(reader io.Reader) RandomBinaryForest {
    // read numTrees first, then read that many trees
    var numTrees int32
    err := binary.Read(reader, binary.LittleEndian, &numTrees)
    if err != nil {
    }

    trees := make([]RandomBinaryTree, numTrees)
    for i := int32(0); i < numTrees; i++ {
        trees[i] = readTreeFromReader(reader)
    }

    return RandomBinaryForest{trees}
}


// TODO: error-handling
func (forest RandomBinaryForest) WriteToWriter(writer io.Writer) {
    err := binary.Write(writer, binary.LittleEndian, int32(len(forest.Trees)))
    if err != nil {
    }
    for _, tree := range forest.Trees {
        tree.writeToWriter(writer)
    }
}
