package main

import (
	"flag"
	"fmt"
	"github.com/moygit/rbf"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"time"
)

const num_features = 784

type Config struct {
	numTrees             int32
	depth                int32
	leafSize             int32
	numFeaturesToCompare int32
	numNeighbors         int32
}

func main() {
	config := getConfig()
	train := readImagesFile("fashion/train_images.bin")
	test := readImagesFile("fashion/test_images.bin")
	trainLabels := readLabelsFile("fashion/train_labels.bin")
	testLabels := readLabelsFile("fashion/test_labels.bin")

	trainStartTime := time.Now().UnixNano()
	forest := rbf.TrainForest(train, config.numTrees, config.depth, config.leafSize, config.numFeaturesToCompare)
	//filename := fmt.Sprintf("mnist_forest_%d_%d_%d_%d.bin", config.numTrees, config.depth, config.leafSize, config.numFeaturesToCompare)
	//f, _ := os.Create(filename)
	//forest.WriteToWriter(f)
	//f.Close()

	//filename := fmt.Sprintf("mnist_forest_%d_%d_%d_%d.bin", config.numTrees, config.depth, config.leafSize, config.numFeaturesToCompare)
	//f, _ := os.Open(filename)
	//forest := rbf.ReadForestFromReader(f)
	//f.Close()

	evalStartTime := time.Now().UnixNano()
	//matchCount := evalL2(forest, train, test, trainLabels, testLabels, config.numNeighbors)
	matchCount := evalPlurality(forest, test, trainLabels, testLabels)
	evalFinishTime := time.Now().UnixNano()

	trainTime := float64(evalStartTime-trainStartTime) / (1000000.0 * 1000.0)
	evalTime := float64(evalFinishTime-evalStartTime) / 1000000.0
	accuracy := float64(matchCount) / float64(len(testLabels))
	fmt.Printf("%d,%d,%d,%.3f,%.5f,%.5f\n", config.numTrees, config.leafSize, config.numFeaturesToCompare, trainTime, evalTime, accuracy)
}

func getConfig() Config {
	numTrees := flag.Int("t", -1, "number of trees")
	depth := flag.Int("d", -1, "depth")
	leafSize := flag.Int("l", -1, "leaf size")
	numFeaturesToCompare := flag.Int("n", -1, "number of features to compare")
	numNeighbors := flag.Int("k", -1, "number of neighbors to consider for plurality-select")
	flag.Parse()
	if *numTrees < 1 || *depth < 1 || *leafSize < 1 || *numFeaturesToCompare < 1 || *numNeighbors < 0 {
		flag.Usage()
		os.Exit(1)
	}
	return Config{int32(*numTrees), int32(*depth), int32(*leafSize), int32(*numFeaturesToCompare), int32(*numNeighbors)}
}

func readImagesFile(filename string) [][]byte {
	file, _ := os.Open(filename)
	defer file.Close()
	buf := make([]byte, getFileSize(file))
	if _, err := io.ReadFull(file, buf); err != nil {
		log.Panicf("error reading %s: %v\n", filename, err)
	}

	buf2d := make([][]byte, len(buf)/num_features)
	for i := range buf2d {
		buf2d[i] = buf[i*num_features : (i+1)*num_features]
	}

	return buf2d
}

func readLabelsFile(filename string) []byte {
	file, _ := os.Open(filename)
	defer file.Close()

	// file format is "label\nlabel\nlabel\n...", so we can just read in the whole file and then drop alternate bytes
	size := getFileSize(file)
	buf := make([]byte, size)
	if _, err := io.ReadFull(file, buf); err != nil {
		log.Panicf("error reading %s: %v\n", filename, err)
	}

	return buf
}

func getFileSize(file *os.File) int64 {
	info, err := file.Stat()
	if err != nil {
		panic(err)
	}
	return info.Size()
}

func evalL2(forest rbf.RandomBinaryForest, train, test [][]byte, trainLabels, testLabels []byte, numNeighbors int32) int {
	matchCount := 0
	for i, image := range test {
		if evalOneImageL2(forest, train, image, trainLabels, testLabels[i], numNeighbors) {
			matchCount += 1
		}
	}
	return matchCount
}

func evalPlurality(forest rbf.RandomBinaryForest, test [][]byte, trainLabels, testLabels []byte) int {
	matchCount := 0
	for i, image := range test {
		if evalOneImagePlurality(forest, image, trainLabels, testLabels[i]) {
			matchCount += 1
		}
	}
	return matchCount
}

func l2Dist(v1, v2 []byte) int {
	distSquared := 0
	for i := range v1 {
		diff := v1[i] - v2[i]
		distSquared += int(diff * diff)
	}
	return distSquared
}

func evalOneImageL2(forest rbf.RandomBinaryForest, train [][]byte, testImage []byte, trainLabels []byte, testLabel byte, numNeighbors int32) bool {
	// query forest
	resultsCount, allTreeResults := forest.FindPointAllResults(testImage)
	// combine all returned results into single slice
	allResults := make([]int32, resultsCount)
	for _, treeResults := range allTreeResults {
		for _, result := range treeResults {
			allResults = append(allResults, result)
		}
	}

	// order them by L2 distance to test point
	sort.Slice(allResults, func(i, j int) bool {
		return l2Dist(testImage, train[allResults[i]]) < l2Dist(testImage, train[allResults[j]])
	})
	// pick out the k nearest ones
	labelCounts := make([]int, 10)
	for _, index := range allResults[:numNeighbors] {
		labelCounts[trainLabels[index]] += 1
	}
	argmaxLabel := argmax(labelCounts)
	return byte(argmaxLabel) == testLabel
}

//func evalOneImageL2(forest rbf.RandomBinaryForest, train [][]byte, testImage []byte, trainLabels []byte, testLabel byte) bool {
//    allResults := forest.FindPointDedupResults(testImage)
//    minDist := math.MaxInt64
//    minDistIndex := int32(-1)
//    for index := range allResults {
//        if dist := l2Dist(testImage, train[index]); dist < minDist {
//            minDist = dist
//            minDistIndex = index
//        }
//    }
//    rbfLabel := trainLabels[minDistIndex]
//    return rbfLabel == testLabel
//}

func evalOneImagePlurality(forest rbf.RandomBinaryForest, image []byte, trainLabels []byte, testLabel byte) bool {
	_, results := forest.FindPointAllResults(image)
	labelCounts := make([]int, 10)
	for _, resultRow := range results {
		for _, result := range resultRow {
			labelCounts[trainLabels[result]] += 1
		}
	}
	argmaxLabel := argmax(labelCounts)
	return byte(argmaxLabel) == testLabel
}

func argmax(values []int) int {
	maxLabelCount := math.MinInt64
	argmaxLabel := -1
	for label := 0; label < len(values); label++ {
		if values[label] > maxLabelCount {
			argmaxLabel = label
			maxLabelCount = values[label]
		}
	}
	return argmaxLabel
}

func printResults(trainStartTime, evalStartTime, evalFinishTime int64, matchCount, totalCount int) {
	fmt.Printf("Total training time: %.2f seconds\n", float64((evalStartTime-trainStartTime)/1000000)/1000.0)
	fmt.Printf("Total eval time: %.2f seconds\n", float64((evalFinishTime-evalStartTime)/1000000)/1000.0)
	fmt.Printf("Average eval time per image: %.3f milliseconds\n", float64((evalFinishTime-evalStartTime)/1000000)/float64(totalCount))
	fmt.Printf("Accuracy: %.2f%%\n", float64(matchCount*100)/float64(totalCount))
}
