package features

import (
	"testing"
)

func TestGetBigramsWithRepeats(t *testing.T) {
	// given:
	b := bigrams{maxBigramCount: 255}
	// when:
	blankBigrams := b.fromString("")
	abcdBigrams := b.fromString("abcdabcdabcd")
	if !testSliceIsSingleValue(blankBigrams, byte(0)) {
		t.Errorf("blank bigrams-with-repeats are wrong :-(")
	}
	if abcdBigrams[0] != 0 ||
		abcdBigrams[1] != 3 || // ab
		abcdBigrams[39] != 3 || // bc
		abcdBigrams[77] != 3 || // cd
		abcdBigrams[111] != 2 || // da
		!testSliceIsSingleValue(abcdBigrams[2:39], byte(0)) || // and there's nothing else
		!testSliceIsSingleValue(abcdBigrams[40:77], byte(0)) ||
		!testSliceIsSingleValue(abcdBigrams[78:111], byte(0)) ||
		!testSliceIsSingleValue(abcdBigrams[112:], byte(0)) {
		t.Errorf("abcdabcdabcd bigrams-with-repeats are wrong :-(")
	}
}

func TestGetBigramsNoRepeats(t *testing.T) {
	// given:
	b := bigrams{maxBigramCount: 1}
	// when:
	blankBigrams := b.fromString("")
	abcdBigrams := b.fromString("abcdabcdabcd")
	if !testSliceIsSingleValue(blankBigrams, byte(0)) {
		t.Errorf("blank bigrams-without-repeats are wrong :-(")
	}
	if abcdBigrams[9] != 0 ||
		abcdBigrams[1] != 1 || // ab
		abcdBigrams[39] != 1 || // bc
		abcdBigrams[77] != 1 || // cd
		abcdBigrams[111] != 1 || // da
		!testSliceIsSingleValue(abcdBigrams[2:39], byte(0)) || // and there's nothing else
		!testSliceIsSingleValue(abcdBigrams[40:77], byte(0)) ||
		!testSliceIsSingleValue(abcdBigrams[78:111], byte(0)) ||
		!testSliceIsSingleValue(abcdBigrams[112:], byte(0)) {
		t.Errorf("abcdabcdabcd bigrams-without-repeats are wrong :-(")
	}
}
