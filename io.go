package rbf


import (
    "bufio"
    "bytes"
    "encoding/binary"
    "encoding/gob"
    "fmt"
    "log"
    "io"
    "os"
    "strings"

    "rbf/features"
)


const READ_CHECK_ERROR_TEMPLATE = "ERROR! Expected %d training strings from training file but read %d\n"


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


func ReadForestFromFile(filename string) RandomBinaryForest {
    // TODO: handle errors
    reader, _ := os.Open(filename)
    defer reader.Close()
    return ReadForestFromReader(reader)
}


// TODO: error-handling
func ReadForestFromReader(reader io.Reader) RandomBinaryForest {
    // read numTrees first, then read that many trees
    readTrees := func() []RandomBinaryTree {
        var numTrees int32
        binary.Read(reader, binary.LittleEndian, &numTrees)
        trees := make([]RandomBinaryTree, numTrees)
        for i := int32(0); i < numTrees; i++ {
            trees[i] = readTreeFromReader(reader)
        }
        return trees
    }

    // read featureSetConfigs
    readFeatureSetConfigs := func() []features.FeatureSetConfig {
        // read count and byte-size
        var numFeatureSetConfigs, gobBytesSize int32
        binary.Read(reader, binary.LittleEndian, &numFeatureSetConfigs)
        binary.Read(reader, binary.LittleEndian, &gobBytesSize)
        // read the bytes
        gobBytes := make([]byte, gobBytesSize)
        binary.Read(reader, binary.LittleEndian, &gobBytes)
        gobReader := bytes.NewBuffer(gobBytes)
        dec := gob.NewDecoder(gobReader)
        // decode the FeatureSetConfigs
        featureSetConfigs := make([]features.FeatureSetConfig, numFeatureSetConfigs)
        for i := range featureSetConfigs {
            // TODO: Is this necessary?
            var fsc features.FeatureSetConfig
            err := dec.Decode(&fsc)
            if err != nil {
                log.Fatal("decode:", err)
            }
            featureSetConfigs[i] = fsc
        }
        return featureSetConfigs
    }

    // TODO: MOVE TO WRAPPER
    // read number of training-strings, then read that many strings
    readTrainingStrings := func() []string {
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

        return trainingStrings
    }

    trees := readTrees()
    featureSetConfigs := readFeatureSetConfigs()
    // TODO: MOVE TO WRAPPER
    trainingStrings := readTrainingStrings()
    calculateFeatures, calculateFeaturesForArray := features.MakeFeatureCalculationFunctions(featureSetConfigs)

    return RandomBinaryForest{trainingStrings, trees, featureSetConfigs, calculateFeatures, calculateFeaturesForArray}
}


func (forest RandomBinaryForest) WriteToWriter(writer io.Writer) {
    // write trees
    writeTrees := func() {
        binary.Write(writer, binary.LittleEndian, int32(len(forest.Trees)))
        for _, tree := range forest.Trees {
            tree.writeToWriter(writer)
        }
    }

    // write featureSetConfigs
    writeFeatureSetConfigs := func() {
        binary.Write(writer, binary.LittleEndian, int32(len(forest.FeatureSetConfigs)))
        var gobWriter bytes.Buffer
        enc := gob.NewEncoder(&gobWriter)
        for i, featureSetConfig := range forest.FeatureSetConfigs {
            err := enc.Encode(&featureSetConfig)
            if err != nil {
                log.Fatalf("unable to encode FeatureSetConfig number %d: %v", i, err)
            }
        }
        gobBytes := gobWriter.Bytes()
        binary.Write(writer, binary.LittleEndian, int32(len(gobBytes)))
        binary.Write(writer, binary.LittleEndian, gobBytes)
    }

    // TODO: MOVE TO WRAPPER
    // write training strings
    writeStrings := func() {
        binary.Write(writer, binary.LittleEndian, int32(len(forest.TrainingStrings)))   // TODO: make this private
        bufWriter := bufio.NewWriter(writer)
        for _, s := range forest.TrainingStrings {   // TODO: make this private
            fmt.Fprintf(writer, "%s\n", s)
        }
        bufWriter.Flush()
    }

    writeTrees()
    writeFeatureSetConfigs()
    // TODO: MOVE TO WRAPPER
    writeStrings()
}

func (forest RandomBinaryForest) WriteToFile(filename string) {
    // TODO: handle errors
    writer, _ := os.Create(filename)
    defer writer.Close()
    forest.WriteToWriter(writer)
}
