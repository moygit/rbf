package features


import (
    "encoding/binary"
    "io"
)


// Feature that gets the last number in a string (mod 256).
// NOTE: We allow the user to specify the number of times they want this repeated.


//----------------------------------------------------------------------------------------------------
// Provide FeatureSetConfig
var DefaultLastNumber LastNumber

type LastNumber struct {
    Count byte
}

func (fn LastNumber) Size() int32 {
    return int32(fn.Count)
}

func (fn LastNumber) FromStringInPlace(input string, featureArray []byte) {
    lastNum := GetLastNumber(input)
    for i := byte(0); i < fn.Count; i++ {
        featureArray[i] = lastNum
    }
}

const last_number_type = int32(12)

func (ln LastNumber) Serialize(writer io.Writer) {
    binary.Write(writer, binary.LittleEndian, last_number_type)
    binary.Write(writer, binary.LittleEndian, int32(ln.Count))
}

func deserialize_last_number(reader io.Reader) FeatureSetConfig {
    var count int32
    binary.Read(reader, binary.LittleEndian, &count)
    return LastNumber{byte(count)}
}
//----------------------------------------------------------------------------------------------------

// Write this from scratch instead of using regular expressions. Much faster.
func GetLastNumber(input string) byte {
    lastCh := byte('-')     // last char before we see the tail of the number
    num := 0
    numPower10 := 0
    for i := len(input) - 1; i >= 0; i-- {
        ch := input[i]
        if ch >= '0' && ch <= '9' {
            if numPower10 > 0 {
                num = (int(ch - '0') * numPower10) + num
                numPower10 = numPower10 * 10
            } else {
                if lastCh < 'a' || lastCh > 'z' {
                    // ok, this is the first (from right) numeric char, and the last char wasn't
                    // alphabetic, so we're now legitimately inside a number
                    num = int(ch - '0')
                    numPower10 = 10
                }
            }
        } else {
            if numPower10 > 0 {
                // we were inside the number but just fell out, so we're done
                return byte(num % 256)
            } else {
                lastCh = ch
            }
        }
    }
    return byte(num % 256)
}

func init() {
    DefaultLastNumber = LastNumber{20}
}
