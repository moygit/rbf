package rbf


import (
    "bufio"
    "encoding/binary"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "strings"
)


const READ_CHECK_ERROR_TEMPLATE = "ERROR! Expected %d training strings from training file but read %d\n"


func readFollowgramsFileIntoSingleArray(filename string) []byte {
    dat, err := ioutil.ReadFile(filename)
    check(err)
    return dat
}

func readFollowgramsFileIntoArrayOfArrays(filename string) [][]byte {
    fullData := readFollowgramsFileIntoSingleArray(filename)
    lines := make([][]byte, len(fullData)/NUM_FEATURES)
    for i := 0; i < len(lines); i++ {
        lines[i] = fullData[(i * NUM_FEATURES):((i + 1) * NUM_FEATURES)]
    }
    return lines
}


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

func ReadTreeFromFile(filename string) RandomBinaryTree {
    // TODO: handle errors
    reader, _ := os.Open(filename)
    defer reader.Close()
    return readTreeFromReader(reader)
}


func (tree RandomBinaryTree) WriteToFile(filename string) {
    // TODO: handle errors
    writer, _ := os.Create(filename)
    defer writer.Close()
    tree.writeToWriter(writer)
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


func ReadForestFromFile(filename string) RandomBinaryForest {
    // TODO: handle errors
    reader, _ := os.Open(filename)
    defer reader.Close()
    return ReadForestFromReader(reader)
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

    // now read training-strings
    var numTrainingStrings int32
    binary.Read(reader, binary.LittleEndian, &numTrainingStrings)
    trainingStrings := make([]string, numTrainingStrings)
    scanner := bufio.NewScanner(reader)
    for i := 0; scanner.Scan(); i++ {
        trainingStrings[i] = strings.TrimSpace(scanner.Text())
    }
 
    // TODO:
    // double-check
    if len(trainingStrings) != int(numTrainingStrings) {
        fmt.Fprintf(os.Stderr, READ_CHECK_ERROR_TEMPLATE, numTrainingStrings, len(trainingStrings))
    }
    return RandomBinaryForest{trainingStrings, trees}
}

func (forest RandomBinaryForest) writeToWriter(writer io.Writer) {
    // write trees
    binary.Write(writer, binary.LittleEndian, int32(len(forest.Trees)))
    for _, tree := range forest.Trees {
        tree.writeToWriter(writer)
    }

    // write training address strings
    binary.Write(writer, binary.LittleEndian, int32(len(forest.trainingStrings)))
    bufWriter := bufio.NewWriter(writer)
    for _, s := range forest.trainingStrings {
        fmt.Fprintf(writer, "%s\n", s)
    }
    bufWriter.Flush()
}

func (forest RandomBinaryForest) WriteToFile(filename string) {
    // TODO: handle errors
    writer, _ := os.Create(filename)
    defer writer.Close()
    forest.writeToWriter(writer)
}
