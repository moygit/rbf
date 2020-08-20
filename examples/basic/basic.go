package main

import (
	"fmt"
	"github.com/moygit/rbf"
)

func main() {
	// Which point is (1,1) closest to, (0,0) or (10,10)?
	points := [][]byte{{0, 0}, {10, 10}}
	queryPoint := []byte{1, 1}

	// Build the forest and query it:
	var numTrees, depth, leafSize, numFeaturesToCompare int32 = 1, 2, 1, 1
	forest := rbf.TrainForest(points, numTrees, depth, leafSize, numFeaturesToCompare)
	count, results := forest.FindPointAllResults(queryPoint)
	nearest := points[results[numTrees-1][count-1]]

	fmt.Printf("Number of results: %d\n", count)                 // 1 point returned
	fmt.Printf("Nearest point to %v: %v\n", queryPoint, nearest) // Nearest point to (1,1) is (0,0)
}
