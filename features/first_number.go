package features

// Feature that gets the first number in a string (mod 256).
// NOTE: We allow the user to specify the number of times they want this feature repeated
// (poor man's weighting).

import (
	"strconv"
	"strings"
)

const first_number_default_count = 20

type firstNumber struct {
	Count byte
}

//----------------------------------------------------------------------------------------------------
// Provide featureSetConfig
func (fn firstNumber) Size() int32 {
	return int32(fn.Count)
}

func (fn firstNumber) FromStringInPlace(input string, featureArray []byte) {
	firstNum := GetFirstNumber(input)
	for i := byte(0); i < fn.Count; i++ {
		featureArray[i] = firstNum
	}
}

func deserializeFirstNumberMap(confMap map[string]string) (config featureSetConfig, ok bool) {
	if countStr, ok := confMap["count"]; ok {
		if count, err := strconv.Atoi(countStr); err == nil {
			return firstNumber{byte(count)}, true
		} else {
			return nil, false
		}
	}
	return firstNumber{byte(first_number_default_count)}, true
}

//----------------------------------------------------------------------------------------------------

// Do this from scratch instead of using regular expressions. Much faster.
func GetFirstNumber(input string) byte {
	num := 0
	inNum := false
	for _, ch := range input {
		if ch >= '0' && ch <= '9' {
			if inNum {
				num = (num * 10) + int(ch-'0')
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
	// if the number was at the end of the string
	return byte(num % 256)
}

// Same as above, but return a string.
func GetFirstNumberAsString(input string) string {
	var num strings.Builder
	inNum := false
	for _, ch := range input {
		if ch >= '0' && ch <= '9' {
			num.WriteByte(byte(ch))
			inNum = true
		} else {
			if inNum {
				if ch >= 'a' && ch <= 'z' {
					// Heuristic: this wasn't an actual number, more like "1st" or "3A"
					num.Reset()
					inNum = false
				} else {
					// we were inside the number but just fell out, so we're done
					return strings.TrimLeft(num.String(), "0")
				}
			}
		}
	}
	// if the number was at the end of the string
	return strings.TrimLeft(num.String(), "0")
}
