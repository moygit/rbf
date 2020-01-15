package features

// Feature that gets the last number in a string (mod 256).
// NOTE: We allow the user to specify the number of times they want this feature repeated
// (poor man's weighting).

import (
	"encoding/binary"
	"io"
	"strings"
)

//----------------------------------------------------------------------------------------------------
// Provide FeatureSetConfig
var DefaultLastNumber LastNumber

type LastNumber struct {
	Count byte
}

func (fn LastNumber) Size() int32 {
	return int32(fn.Count)
}

func (fn LastNumber) fromStringInPlace(input string, featureArray []byte) {
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

// Do this from scratch instead of using regular expressions. Much faster.
func GetLastNumber(input string) byte {
	lastCh := byte('-') // keep track of the last char we saw as we scan from the right before we saw the tail of the number
	num := 0
	numPower10 := 0 // use this to track which digit we're on (and hence also whether we're inside a number)
	for i := len(input) - 1; i >= 0; i-- {
		ch := input[i]
		if ch >= '0' && ch <= '9' {
			if numPower10 > 0 {
				// we were already inside a number
				num = (int(ch-'0') * numPower10) + num
				numPower10 = numPower10 * 10
			} else {
				if lastCh < 'a' || lastCh > 'z' {
					// this is the first (from right) numeric char, and the last char wasn't alphabetic
					// (i.e. it was punctuation or whitespace), so we're now legitimately inside a number
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
	// if the number was at the (left) end of the string
	return byte(num % 256)
}

// Same as above, but return a string.
func GetLastNumberAsString(input string) string {
	reverse := func(revDigits []byte) string {
		var num strings.Builder
		for i := len(revDigits) - 1; i >= 0; i-- {
			num.WriteByte(revDigits[i])
		}
		return num.String()
	}

	lastCh := byte('-') // keep track of the last char we saw as we scan from the right before we saw the tail of the number
	num := make([]byte, 0)
	inNum := false
	for i := len(input) - 1; i >= 0; i-- {
		ch := input[i]
		if ch >= '0' && ch <= '9' {
			if inNum {
				// we were already inside a number
				num = append(num, byte(ch))
			} else {
				if lastCh < 'a' || lastCh > 'z' {
					// this is the first (from right) numeric char, and the last char wasn't alphabetic
					// (i.e. it was punctuation or whitespace), so we're now legitimately inside a number
					num = append(num, byte(ch))
					inNum = true
				}
			}
		} else {
			if inNum {
				// we were inside the number but just fell out, so we're done
				return reverse(num)
			} else {
				lastCh = ch
			}
		}
	}
	// if the number was at the (left) end of the string
	return reverse(num)
}

func init() {
	DefaultLastNumber = LastNumber{20}
}
