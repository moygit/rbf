package features


import (
    "testing"
)


func TestFeatureSetConfig(t *testing.T) {
    sliceIsSingleValue := func (slice []byte, val byte) bool {
        for i := 0; i < len(slice); i++ {
            if slice[i] != val {
                return false
            }
        }
        return true
    }

    // given/when:
    featureSetConfig := []FeatureSetConfig{ Followgrams{3}, Followgrams{3} }
    calculateFeatures, _ := MakeFeatureCalculationFunctions(featureSetConfig)
    followgrams := calculateFeatures("abcdefgh")

    // then:
    if !sliceIsSingleValue(followgrams[1:4], byte(1)) ||        // a
       !sliceIsSingleValue(followgrams[39:42], byte(1)) ||      // b
       !sliceIsSingleValue(followgrams[77:80], byte(1)) ||      // c
       !sliceIsSingleValue(followgrams[115:118], byte(1)) ||    // d
       !sliceIsSingleValue(followgrams[153:156], byte(1)) ||    // e
       !sliceIsSingleValue(followgrams[191:193], byte(1)) ||    // f (only 2)
       !sliceIsSingleValue(followgrams[229:230], byte(1)) ||    // g (only 1)
       !sliceIsSingleValue(followgrams[num_followgrams + 1:num_followgrams + 4], byte(1)) ||        // a
       !sliceIsSingleValue(followgrams[num_followgrams + 39:num_followgrams + 42], byte(1)) ||      // b
       !sliceIsSingleValue(followgrams[num_followgrams + 77:num_followgrams + 80], byte(1)) ||      // c
       !sliceIsSingleValue(followgrams[num_followgrams + 115:num_followgrams + 118], byte(1)) ||    // d
       !sliceIsSingleValue(followgrams[num_followgrams + 153:num_followgrams + 156], byte(1)) ||    // e
       !sliceIsSingleValue(followgrams[num_followgrams + 191:num_followgrams + 193], byte(1)) ||    // f (only 2)
       !sliceIsSingleValue(followgrams[num_followgrams + 229:num_followgrams + 230], byte(1)) ||    // g (only 1)
       false {
        t.Errorf("abcdefgh 3-followgrams are wrong :-(")
    }
    if !sliceIsSingleValue(followgrams[0:1], byte(0)) ||        // before a
       !sliceIsSingleValue(followgrams[4:39], byte(0)) ||       // between a and b
       !sliceIsSingleValue(followgrams[42:77], byte(0)) ||      // between b and c
       !sliceIsSingleValue(followgrams[80:115], byte(0)) ||     // between c and d
       !sliceIsSingleValue(followgrams[118:153], byte(0)) ||    // between d and e
       !sliceIsSingleValue(followgrams[156:191], byte(0)) ||    // between e and f
       !sliceIsSingleValue(followgrams[193:229], byte(0)) ||    // between f and g
       !sliceIsSingleValue(followgrams[230:num_followgrams], byte(0)) ||       // after g
       !sliceIsSingleValue(followgrams[num_followgrams+0:num_followgrams + 1], byte(0)) ||        // before a
       !sliceIsSingleValue(followgrams[num_followgrams+4:num_followgrams + 39], byte(0)) ||       // between a and b
       !sliceIsSingleValue(followgrams[num_followgrams+42:num_followgrams + 77], byte(0)) ||      // between b and c
       !sliceIsSingleValue(followgrams[num_followgrams+80:num_followgrams + 115], byte(0)) ||     // between c and d
       !sliceIsSingleValue(followgrams[num_followgrams+118:num_followgrams + 153], byte(0)) ||    // between d and e
       !sliceIsSingleValue(followgrams[num_followgrams+156:num_followgrams + 191], byte(0)) ||    // between e and f
       !sliceIsSingleValue(followgrams[num_followgrams+193:num_followgrams + 229], byte(0)) ||    // between f and g
       !sliceIsSingleValue(followgrams[num_followgrams+230:], byte(0)) ||       // after g
       false {
        t.Errorf("abcdefgh 3-followgrams are wrong :-(")
    }

    // given/when:
    featureSetConfig = []FeatureSetConfig{ Followgrams{6}, Followgrams{6} }
    calculateFeatures, _ = MakeFeatureCalculationFunctions(featureSetConfig)
    followgrams = calculateFeatures("aaaaaaaa")

    // then:
    if followgrams[0] != 27 || followgrams[num_followgrams] != 27 {
        t.Errorf("aa count is %d; expected %d", followgrams[0], 27)
    }
    if !sliceIsSingleValue(followgrams[1:num_followgrams], byte(0)) ||
       !sliceIsSingleValue(followgrams[num_followgrams+1:], byte(0)) {
        t.Errorf("got non-zero count for some non-aa value")
    }
}
