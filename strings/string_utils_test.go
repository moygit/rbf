package strings


import (
    "math"
    "testing"

    "fmt"
    "github.com/adrg/strutil"
    "github.com/adrg/strutil/metrics"
)


const FLOAT_TOLERANCE = 0.000001

const s0 = "1 r de richelieu paris 75010"
const s1 = "1 r xdes xrichelieux xparisx 75010"  // typos
const s2 = "1 r de richelieu 75010 paris"    // just a postal-code switch, should still match
const s3 = "e r lr sdi a1eihp0 c0175ieur"    // permutation of s0, don't want this to match

const t0 = "06 louxembourg paris 75006 fr"
const t1 = "96 rue beaubourg paris 75003 fr"


func TestJaccardSimilarity(t *testing.T) {
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
}


func TestOneDirectionalJaccard(t *testing.T) {
    fmt.Printf("OneDirectionalJaccard(\"%s\", \"%s\") == %f\n", s0, s1, OneDirectionalJaccard2GramSimilarity(s0, s1))
    fmt.Printf("OneDirectionalJaccard(\"%s\", \"%s\") == %f\n", s0, s2, OneDirectionalJaccard2GramSimilarity(s0, s2))
    fmt.Printf("OneDirectionalJaccard(\"%s\", \"%s\") == %f\n", s0, s3, OneDirectionalJaccard2GramSimilarity(s0, s3))

    fmt.Printf("OneDirectionalJaccard(\"%s\", \"%s\") == %f\n", t0, t1, OneDirectionalJaccard2GramSimilarity(t0, t1))
    fmt.Println()
}


func TestLevenshtein(t *testing.T) {
    lev := metrics.NewLevenshtein()
    fmt.Printf("Levenshtein(\"%s\", \"%s\") == %f\n", s0, s1, strutil.Similarity(s0, s1, lev))
    fmt.Printf("Levenshtein(\"%s\", \"%s\") == %f\n", s0, s2, strutil.Similarity(s0, s2, lev))
    fmt.Printf("Levenshtein(\"%s\", \"%s\") == %f\n", s0, s3, strutil.Similarity(s0, s3, lev))

    fmt.Printf("Levenshtein(\"%s\", \"%s\") == %f\n", t0, t1, strutil.Similarity(t0, t1, lev))
    fmt.Println()
}

func TestSmithWatermanAndSorensenDice(t *testing.T) {
//    swg := metrics.NewSmithWatermanGotoh()
//    fmt.Printf("SmithWatermanGotoh(\"%s\", \"%s\") == %f\n", s0, s1, strutil.Similarity(s0, s1, swg))
//    fmt.Printf("SmithWatermanGotoh(\"%s\", \"%s\") == %f\n", s0, s2, strutil.Similarity(s0, s2, swg))
//    fmt.Printf("SmithWatermanGotoh(\"%s\", \"%s\") == %f\n", s0, s3, strutil.Similarity(s0, s3, swg))
//    fmt.Println()
// 
//    sd := metrics.NewSorensenDice()
//    sd.NgramSize = 2
//    fmt.Printf("SorensenDice(\"%s\", \"%s\") == %f\n", s0, s1, strutil.Similarity(s0, s1, sd))
//    fmt.Printf("SorensenDice(\"%s\", \"%s\") == %f\n", s0, s2, strutil.Similarity(s0, s2, sd))
//    fmt.Printf("SorensenDice(\"%s\", \"%s\") == %f\n", s0, s3, strutil.Similarity(s0, s3, sd))
//    fmt.Println()
//    // fmt.Printf("%.2f\n", similarity) // Output: 0.82
// 
//    sd.NgramSize = 2
//    x1 := "model town near reliance wala store house no 18r model town ludhiana punjab 141002 in"
//    x2 := "model town near hiazr maker saloon kothi no 18r model townnear sikand car ludhiana punjab 141002 in"
//    x3 := "model town near hiazr maker saloon kothi no 18r opp prefcet car care house no 18r model town ludhiana punjab 141002 in"
//    x4 := "model town near hiazr mkaer salon kothi no 18r near sikand cvar world ludhiana punjab 141002 in"
//    x5 := "model town near hiazr mkaer saloon kothi no 18r modle town near sikand cra ludhiana punjab 141002 in"
//    x6 := "model town near hiazr mkaer saooloon kothi no 18r model town near dugri wali road ludhiana punjab 141002 in"
//    x7 := "model town near jolly dep store dugri rpoad houseno 18r kothi no 18r model town ludhiana punjab 141002 in"
//    fmt.Printf("SmithWatermanGotoh(\"%s\", \"%s\") == %f\n", x2, x1, strutil.Similarity(x2, x1, swg))
//    fmt.Printf("SmithWatermanGotoh(\"%s\", \"%s\") == %f\n", x2, x3, strutil.Similarity(x2, x3, swg))
//    fmt.Printf("SmithWatermanGotoh(\"%s\", \"%s\") == %f\n", x2, x4, strutil.Similarity(x2, x4, swg))
//    fmt.Printf("SmithWatermanGotoh(\"%s\", \"%s\") == %f\n", x2, x5, strutil.Similarity(x2, x5, swg))
//    fmt.Printf("SmithWatermanGotoh(\"%s\", \"%s\") == %f\n", x2, x6, strutil.Similarity(x2, x6, swg))
//    fmt.Printf("SmithWatermanGotoh(\"%s\", \"%s\") == %f\n", x2, x7, strutil.Similarity(x2, x7, swg))
//    fmt.Println()
//    fmt.Printf("SorensenDice(\"%s\", \"%s\") == %f\n", x2, x1, strutil.Similarity(x2, x1, sd))
//    fmt.Printf("SorensenDice(\"%s\", \"%s\") == %f\n", x2, x3, strutil.Similarity(x2, x3, sd))
//    fmt.Printf("SorensenDice(\"%s\", \"%s\") == %f\n", x2, x4, strutil.Similarity(x2, x4, sd))
//    fmt.Printf("SorensenDice(\"%s\", \"%s\") == %f\n", x2, x5, strutil.Similarity(x2, x5, sd))
//    fmt.Printf("SorensenDice(\"%s\", \"%s\") == %f\n", x2, x6, strutil.Similarity(x2, x6, sd))
//    fmt.Printf("SorensenDice(\"%s\", \"%s\") == %f\n", x2, x7, strutil.Similarity(x2, x7, sd))
}
