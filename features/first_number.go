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
    firstNum := GetFirstNumber(input)
    for i := byte(0); i < fn.Count; i++ {
        featureArray[i] = firstNum
    }
}
//----------------------------------------------------------------------------------------------------

// Write this from scratch instead of using regular expressions. Much faster.
func GetFirstNumber(input string) byte {
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
                if ch >= 'a' && ch <= 'z' {
                    // Heuristic: this wasn't an actual number, more like "1st" or "3A"
                    num = 0
                    inNum = false
                } else {
                    // we were inside the number but just fell out, so we're done
                    return byte(num % 256)
                }
            }
        }
    }
    return byte(num % 256)
}

func init() {
    gob.Register(FirstNumber{})
    DefaultFirstNumber = FirstNumber{20}
}
