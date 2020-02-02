package features

import "testing"

func TestGetFollowgrams(t *testing.T) {
	// given:
	f := followgrams{3}
	// when:
	grams := f.fromString("abcdefgh")
	// then:
	if !testSliceIsSingleValue(grams[1:4], byte(1)) || // a
		!testSliceIsSingleValue(grams[39:42], byte(1)) || // b
		!testSliceIsSingleValue(grams[77:80], byte(1)) || // c
		!testSliceIsSingleValue(grams[115:118], byte(1)) || // d
		!testSliceIsSingleValue(grams[153:156], byte(1)) || // e
		!testSliceIsSingleValue(grams[191:193], byte(1)) || // f (only 2)
		!testSliceIsSingleValue(grams[229:230], byte(1)) || // g (only 1)
		false {
		t.Errorf("abcdefgh 3-followgrams are wrong :-(")
	}
	if !testSliceIsSingleValue(grams[0:1], byte(0)) || // before a
		!testSliceIsSingleValue(grams[4:39], byte(0)) || // between a and b
		!testSliceIsSingleValue(grams[42:77], byte(0)) || // between b and c
		!testSliceIsSingleValue(grams[80:115], byte(0)) || // between c and d
		!testSliceIsSingleValue(grams[118:153], byte(0)) || // between d and e
		!testSliceIsSingleValue(grams[156:191], byte(0)) || // between e and f
		!testSliceIsSingleValue(grams[193:229], byte(0)) || // between f and g
		!testSliceIsSingleValue(grams[230:], byte(0)) || // after g
		false {
		t.Errorf("abcdefgh 3-followgrams are wrong :-(")
	}

	// given:
	f = followgrams{6}
	// when:
	grams = f.fromString("aaaaaaaa")
	// then:
	if grams[0] != 27 {
		t.Errorf("aa count is %d; expected %d", grams[0], 27)
	}
	if !testSliceIsSingleValue(grams[1:], byte(0)) {
		t.Errorf("got non-zero count for some non-aa value")
	}
}
