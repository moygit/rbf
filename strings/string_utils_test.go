package strings


import (
    "math"
    "testing"

    "fmt"
    "github.com/adrg/strutil"
    "github.com/adrg/strutil/metrics"
)


const FLOAT_TOLERANCE = 0.000001


func TestGetFollowgrams(t *testing.T) {
    sliceIsSingleValue := func (slice []byte, val byte) bool {
        for i := 0; i < len(slice); i++ {
            if slice[i] != val {
                return false
            }
        }
        return true
    }

    // given/when:
    followgrams := GetFollowgramsWithWindowSize("abcdefgh", 3)
    if !sliceIsSingleValue(followgrams[1:4], byte(1)) ||        // a
       !sliceIsSingleValue(followgrams[39:42], byte(1)) ||      // b
       !sliceIsSingleValue(followgrams[77:80], byte(1)) ||      // c
       !sliceIsSingleValue(followgrams[115:118], byte(1)) ||    // d
       !sliceIsSingleValue(followgrams[153:156], byte(1)) ||    // e
       !sliceIsSingleValue(followgrams[191:193], byte(1)) ||    // f (only 2)
       !sliceIsSingleValue(followgrams[229:230], byte(1)) ||    // g (only 1)
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
       !sliceIsSingleValue(followgrams[230:], byte(0)) ||       // after g
       false {
        t.Errorf("abcdefgh 3-followgrams are wrong :-(")
    }
    if !sliceIsSingleValue(followgrams[0:1], byte(0)) ||
       !sliceIsSingleValue(followgrams[1:4], byte(1)) ||
       !sliceIsSingleValue(followgrams[4:39], byte(0)) ||
       !sliceIsSingleValue(followgrams[39:42], byte(1)) ||
       !sliceIsSingleValue(followgrams[42:76], byte(0)) {
        t.Errorf("abcdefgh 3-followgrams are wrong :-(")
    }

    // given/when:
    followgrams = GetFollowgramsWithWindowSize("aaaaaaaa", 6)
    if followgrams[0] != 27 {
        t.Errorf("aa count is %d; expected %d", followgrams[0], 27)
    }
    if !sliceIsSingleValue(followgrams[1:], byte(0)) {
        t.Errorf("got non-zero count for some non-aa value")
    }
}


func TestJaccardSimilarity(t *testing.T) {
    s0 := "1 r de richelieu paris 75010"
    s1 := "1 r xdes xrichelieux xparisx 75010"  // typos
    s2 := "1 r de richelieu 75010 paris"    // just a postal-code switch, should still match
    s3 := "e r lr sdi a1eihp0 c0175ieur"    // permutation of s0, don't want this to match

    testOneCase := func(s1, s2 string, simFunc func(string, string) float64, funcName string, expSim float64) {
        sim := simFunc(s1, s2)
        if !(math.Abs(expSim - sim) < FLOAT_TOLERANCE) {
            t.Errorf("%s(%s, %s) == %f; expected %f", funcName, s1, s2, sim, expSim)
        }
    }

    testOneCase(s0, s1, Jaccard2gramSimilarity, "Jaccard2gramSimilarity", 0.578947)
    testOneCase(s0, s2, Jaccard2gramSimilarity, "Jaccard2gramSimilarity", 0.928571)
    testOneCase(s0, s3, Jaccard2gramSimilarity, "Jaccard2gramSimilarity", 0.148936)

    testOneCase(s0, s1, JaccardFollowgramSimilarity, "JaccardFollowgramSimilarity", 0.673797)
    testOneCase(s0, s2, JaccardFollowgramSimilarity, "JaccardFollowgramSimilarity", 0.830508)
    testOneCase(s0, s3, JaccardFollowgramSimilarity, "JaccardFollowgramSimilarity", 0.632829)

    // similarity := strutil.Similarity("graph", "giraffe", metrics.NewLevenshtein())
    // fmt.Printf("%.2f\n", similarity) // Output: 0.43

    // similarity = strutil.Similarity("think", "tank", metrics.NewJaro())
    // fmt.Printf("%.2f\n", similarity) // Output: 0.78

    lev := metrics.NewLevenshtein()
    fmt.Printf("Levenshtein(\"%s\", \"%s\") == %f\n", s0, s1, strutil.Similarity(s0, s1, lev))
    fmt.Printf("Levenshtein(\"%s\", \"%s\") == %f\n", s0, s2, strutil.Similarity(s0, s2, lev))
    fmt.Printf("Levenshtein(\"%s\", \"%s\") == %f\n", s0, s3, strutil.Similarity(s0, s3, lev))
    fmt.Println()

    swg := metrics.NewSmithWatermanGotoh()
    fmt.Printf("SmithWatermanGotoh(\"%s\", \"%s\") == %f\n", s0, s1, strutil.Similarity(s0, s1, swg))
    fmt.Printf("SmithWatermanGotoh(\"%s\", \"%s\") == %f\n", s0, s2, strutil.Similarity(s0, s2, swg))
    fmt.Printf("SmithWatermanGotoh(\"%s\", \"%s\") == %f\n", s0, s3, strutil.Similarity(s0, s3, swg))
    fmt.Println()

    sd := metrics.NewSorensenDice()
    sd.NgramSize = 2
    fmt.Printf("SorensenDice(\"%s\", \"%s\") == %f\n", s0, s1, strutil.Similarity(s0, s1, sd))
    fmt.Printf("SorensenDice(\"%s\", \"%s\") == %f\n", s0, s2, strutil.Similarity(s0, s2, sd))
    fmt.Printf("SorensenDice(\"%s\", \"%s\") == %f\n", s0, s3, strutil.Similarity(s0, s3, sd))
    fmt.Println()
    // fmt.Printf("%.2f\n", similarity) // Output: 0.82

    sd.NgramSize = 2
    x1 := "model town near reliance wala store house no 18r model town ludhiana punjab 141002 in"
    x2 := "model town near hiazr maker saloon kothi no 18r model townnear sikand car ludhiana punjab 141002 in"
    x3 := "model town near hiazr maker saloon kothi no 18r opp prefcet car care house no 18r model town ludhiana punjab 141002 in"
    x4 := "model town near hiazr mkaer salon kothi no 18r near sikand cvar world ludhiana punjab 141002 in"
    x5 := "model town near hiazr mkaer saloon kothi no 18r modle town near sikand cra ludhiana punjab 141002 in"
    x6 := "model town near hiazr mkaer saooloon kothi no 18r model town near dugri wali road ludhiana punjab 141002 in"
    x7 := "model town near jolly dep store dugri rpoad houseno 18r kothi no 18r model town ludhiana punjab 141002 in"
    fmt.Printf("SmithWatermanGotoh(\"%s\", \"%s\") == %f\n", x2, x1, strutil.Similarity(x2, x1, swg))
    fmt.Printf("SmithWatermanGotoh(\"%s\", \"%s\") == %f\n", x2, x3, strutil.Similarity(x2, x3, swg))
    fmt.Printf("SmithWatermanGotoh(\"%s\", \"%s\") == %f\n", x2, x4, strutil.Similarity(x2, x4, swg))
    fmt.Printf("SmithWatermanGotoh(\"%s\", \"%s\") == %f\n", x2, x5, strutil.Similarity(x2, x5, swg))
    fmt.Printf("SmithWatermanGotoh(\"%s\", \"%s\") == %f\n", x2, x6, strutil.Similarity(x2, x6, swg))
    fmt.Printf("SmithWatermanGotoh(\"%s\", \"%s\") == %f\n", x2, x7, strutil.Similarity(x2, x7, swg))
    fmt.Println()
    fmt.Printf("SorensenDice(\"%s\", \"%s\") == %f\n", x2, x1, strutil.Similarity(x2, x1, sd))
    fmt.Printf("SorensenDice(\"%s\", \"%s\") == %f\n", x2, x3, strutil.Similarity(x2, x3, sd))
    fmt.Printf("SorensenDice(\"%s\", \"%s\") == %f\n", x2, x4, strutil.Similarity(x2, x4, sd))
    fmt.Printf("SorensenDice(\"%s\", \"%s\") == %f\n", x2, x5, strutil.Similarity(x2, x5, sd))
    fmt.Printf("SorensenDice(\"%s\", \"%s\") == %f\n", x2, x6, strutil.Similarity(x2, x6, sd))
    fmt.Printf("SorensenDice(\"%s\", \"%s\") == %f\n", x2, x7, strutil.Similarity(x2, x7, sd))
}
