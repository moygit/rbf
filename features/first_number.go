package features


// Feature that gets the first number in a string (mod 256).
// NOTE: We allow the user to specify the number of times they want this repeated.


import "encoding/gob"


//----------------------------------------------------------------------------------------------------
// Wrapper around FirstNumber to provide FeatureSetConfig
var DefaultFirstNumber FirstNumber

type FirstNumber struct {
    Count byte
}

func (fn FirstNumber) Size() int32 {
    return int32(fn.Count)
}

func (fn FirstNumber) FromStringInPlace(input string, featureArray []byte) {
    firstNum := getFirstNumber(input)
    for i := byte(0); i < fn.Count; i++ {
        featureArray[i] = firstNum
    }
}
//----------------------------------------------------------------------------------------------------

func getFirstNumber(input string) byte {
    num := 0
    inNum := false
    for _, ch := range input {
        if ch >= '0' && ch <= '9' {
            if inNum {
                num = (num * 10) + int(ch - '0')
            } else {
                num = int(ch - '0')
                inNum = true
            }
        } else {
            if inNum {
                // we were inside the number but just fell out, so we're done
                return byte(num % 256)
            }
        }
    }
    return byte(num % 256)
}

func init() {
    gob.Register(OccurrenceCounts{})
    DefaultFirstNumber = FirstNumber{10}
}
