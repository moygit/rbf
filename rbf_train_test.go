package rbf

import (
	"reflect"
	"testing"
)

var result int32

func BenchmarkGetSingleFeatureFrequencies(b *testing.B) {
	// run getSingleFeatureFrequencies b.N times
	for n := 0; n < b.N; n++ {
		localFeatureArray := [][]byte{{0}, {0}, {5}, {5}, {5}, {5}, {7}, {7}, {7}, {7}}
		localRowIndex := []int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		featureNum := int32(0)
		indexStart := int32(0)
		indexEnd := int32(10)
		freqs, weightedTotal := getSingleFeatureFrequencies(localRowIndex, localFeatureArray, featureNum, indexStart, indexEnd)

		result += weightedTotal + freqs[0]
	}
}

func TestGetSingleFeatureFrequencies(t *testing.T) {
	// given
	localFeatureArray := [][]byte{{0}, {0}, {5}, {5}, {5}, {5}, {7}, {7}, {7}, {7}}
	localRowIndex := []int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	featureNum := int32(0)
	indexStart := int32(0)
	indexEnd := int32(10)
	// when
	freqs, weightedTotal := getSingleFeatureFrequencies(localRowIndex, localFeatureArray, featureNum, indexStart, indexEnd)
	// then
	expectedInit := []int32{2, 0, 0, 0, 0, 4, 0, 4} // 2 0's, 4 5's, 4 7's,...
	if !reflect.DeepEqual(freqs[:8], expectedInit) {
		t.Errorf("getSingleFeatureFrequencies initial is %v; expected %v", freqs[:8], expectedInit)
	}
	expectedTail := make([]int32, 256-8) // ...and nothing else
	if !reflect.DeepEqual(freqs[8:], expectedTail) {
		t.Errorf("getSingleFeatureFrequencies tail is %v; expected %v", freqs[8:], expectedTail)
	}
	if weightedTotal != 48 {
		t.Errorf("weightedTotal == %d; expected %d", weightedTotal, 48)
	}
}

func TestSplitOneFeature(t *testing.T) {
	runOneTest := func(L []int32, expTotalMoment float32, expPos int32, expLeftCount int32) {
		var moment, count int32
		for i, x := range L {
			moment += int32(i) * int32(x)
			count += int32(x)
		}
		// when
		totalMoment, pos, leftCount := splitOneFeature(L, moment, count)
		// then
		if (totalMoment != expTotalMoment) || (pos != expPos) || (leftCount != expLeftCount) {
			t.Errorf("(totalMoment, pos, leftCount) == (%f, %d, %d); expected (%f, %d, %d)",
				totalMoment, pos, leftCount, expTotalMoment, expPos, expLeftCount)
		}
	}

	L1 := []int32{10, 5, 4, 0, 0, 11, 12, 13}
	runOneTest(L1, 122.5, 5, 30)

	L2 := []int32{10, 0, 0, 0, 0}
	runOneTest(L2, 5.0, 0, 10)

	L3 := []int32{1, 1, 1, 1, 1}
	runOneTest(L3, 6.5, 2, 3)
}

func TestGetBestFeature(t *testing.T) {
	// given:
	featureFrequencies := [][]int32{{1, 1, 1, 1, 1}, {5, 0, 0, 0, 0}}
	weightedTotals := []int32{10, 0} // 10 == (1 * 0) + (1 * 1) + (1 * 2) + (1 * 3) + (1 * 4) + (1 * 5)
	//  0 == (5 * 0) + (0 * 1) + ... + (0 * 4)
	totalCount := int32(5)
	// when:
	bestFeatureNum, splitValue := getBestFeature(featureFrequencies, weightedTotals, totalCount)
	// then f0 is "better" than f1 because its variance is higher
	expBestFeature, expSplitValue := int32(0), byte(2)
	if (bestFeatureNum != expBestFeature) || (splitValue != expSplitValue) {
		t.Errorf("(bestFeatureNum, splitValue) == (%d, %d); expected (%d, %d)\n",
			bestFeatureNum, splitValue, expBestFeature, expSplitValue)
	}
}

func TestGetSimpleBestFeature(t *testing.T) {
	// given:
	featureFrequencies := [][]int32{{1, 1, 1, 1, 1}, {5, 0, 0, 0, 0}}
	weightedTotals := []int32{10, 0} // 10 == (1 * 0) + (1 * 1) + (1 * 2) + (1 * 3) + (1 * 4) + (1 * 5)
	//  0 == (5 * 0) + (0 * 1) + ... + (0 * 4)
	totalCount := int32(5)
	// when:
	bestFeatureNum, splitValue := getSimpleBestFeature(featureFrequencies, weightedTotals, totalCount)
	// then f0 is "better" than f1 because its split is closer to the median
	expBestFeature, expSplitValue := int32(0), byte(2)
	if (bestFeatureNum != expBestFeature) || (splitValue != expSplitValue) {
		t.Errorf("(bestFeatureNum, splitValue) == (%d, %d); expected (%d, %d)\n",
			bestFeatureNum, splitValue, expBestFeature, expSplitValue)
	}
}

func TestQuickPartition(t *testing.T) {
	// setup boilerplate (only 1 feature, we'll sort a 6-length array):
	featureNum, indexStart, indexEnd := int32(0), int32(0), int32(6)

	runOneTest := func(rowIndex []int32, features [][]byte, splitValue byte, expSplit int32, expRowIndex []int32) {
		split := quickPartition(rowIndex, features, indexStart, indexEnd, featureNum, splitValue)
		if split != expSplit {
			t.Errorf("split == %d; expected %d", split, expSplit)
		}
		if !reflect.DeepEqual(rowIndex, expRowIndex) {
			t.Errorf("rowIndex == %v; expected %v", rowIndex, expRowIndex)
		}
	}

	// given a reverse-sorted list:
	rowIndex := []int32{0, 1, 2, 3, 4, 5}
	features := [][]byte{{15}, {14}, {13}, {12}, {11}, {10}}
	// when we split at 12 (largest integer <= median):
	var splitValue byte = 12
	// then we expect to split at position 3 to get feature order [10, 11, 12   |   13, 14, 15]
	expSplit := int32(3)
	expRowIndex := []int32{5, 4, 3, 2, 1, 0}
	runOneTest(rowIndex, features, splitValue, expSplit, expRowIndex)

	// given a particular (random) order:
	rowIndex = []int32{0, 1, 2, 3, 4, 5}
	features = [][]byte{{11}, {10}, {14}, {12}, {15}, {13}}
	// when we split at 12 (largest integer <= median):
	splitValue = 12
	// then we expect to split at position 3 to get feature order [11, 10, 12   |   14, 15, 13]
	expSplit = 3
	expRowIndex = []int32{0, 1, 3, 2, 4, 5}
	runOneTest(rowIndex, features, splitValue, expSplit, expRowIndex)

	// given a particular (random) order:
	rowIndex = []int32{0, 1, 2, 3, 4, 5}
	features = [][]byte{{11}, {10}, {14}, {12}, {15}, {13}}
	// when split-value is less than all the values in the array
	splitValue = 2
	// then there's no split
	expSplit = 0                            // there's no split
	expRowIndex = []int32{0, 1, 2, 3, 4, 5} // and nothing gets moved
	runOneTest(rowIndex, features, splitValue, expSplit, expRowIndex)

	// given empty lists of rows and features:
	rowIndex = []int32{}
	features = [][]byte{}
	indexStart, indexEnd = 0, 0
	// for some split-value
	splitValue = 10
	// "Dude, there's nothing to split here!"
	expSplit = 0
	expRowIndex = []int32{}
	runOneTest(rowIndex, features, splitValue, expSplit, expRowIndex)
}
