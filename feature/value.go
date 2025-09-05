package feature

import (
	"errors"
	"fmt"
	"strings"
)

// Default attribute values
const (
	D_RAWSTR   string = "NO_VALUE"
	D_ISBOOL   bool   = false
	D_HASQUOTE bool   = true
)

// Value structure
type Value struct {
	RawStr   string
	IsBool   bool
	HasQuote bool
}

// Value constructor
func NewValue(s ...string) (Value, error) {
	v := Value{D_RAWSTR, D_ISBOOL, D_HASQUOTE}
	if s == nil {
		// Calling the constructor without string value yields
		// to a Value object with default attribute values
		return v, nil
	} else {
		if len(s) > 1 {
			// Only one string should be given
			return v, errors.New("when creating Value object, you are supposed to provide only one string")
		}

		// If the provided string is empty then consider it is
		// a boolean objet
		if s[0] == "" {
			v.IsBool = true
			return v, nil
		}

		// Else, check if the provided string start with quote character
		str, pre := strings.CutPrefix(s[0], "\"") // Keep the original string unchanged
		str, suf := strings.CutSuffix(str, "\"")
		v.RawStr = str
		// It has a starting quote
		if pre {
			// Then it must ends with a quote character
			if suf {
				v.HasQuote = true
			} else {
				// Missing terminal quote
				return v, fmt.Errorf("the provided string value miss a terminal quote (%s)", s[0])
			}
		} else {
			if suf {
				// Missing starting quote
				return v, fmt.Errorf("the provided string value miss a starting quote (%s)", s[0])
			}
		}

		// Return the initialized Value object
		return v, nil
	}
}

// Get a formatted string of the stored value
func (v *Value) ToString() string {
	// If boolean and/or RawStr is equal to D_RAWSTR then return an empty string
	if v.IsBool || v.RawStr == D_RAWSTR {
		return ""
	}

	// Add quote if necessary
	if v.HasQuote {
		return fmt.Sprintf("\"%s\"", v.RawStr)
	}

	// No quote required
	return v.RawStr
}

/* Old version of ToString method, to be spread into qualifier class...
func (v *Value) ToString(tag, prefix string, nchar int) string {
	// Check if the nchar is not too small
	lenTag := len(tag)
	lenPre := len(prefix)
	if nchar <= lenTag+lenPre+4 {
		panic("The number of character per line is not sufficient to write the value.")
	}

	// Simple case, the value is boolean
	if v.IsBool {
		return fmt.Sprintf("%s%s%s", prefix, D_TAGPREFIX, tag)
	} else {
		// Prepare string value to return
		tmpStr := v.RawStr

		// Check if quotes are required
		if v.HasQuote {
			tmpStr = "\"" + v.RawStr + "\""
		}

		// Add the qualifier tag (if necessary)
		if lenTag > 0 {
			tmpStr = D_TAGPREFIX + tag + "=" + tmpStr
		}

		// Check if the value is small enough to be printed in a single line
		lenStr := len(tmpStr)
		lenRem := nchar - lenPre // Remaining char excepting the prefix
		if lenStr <= lenRem {
			return fmt.Sprintf("%s%s", prefix, tmpStr)
		} else {
			// Split value into multiple lines
			reSpace := regexp.MustCompile(`\s+`)
			reComa := regexp.MustCompile(`\,`)
			var subValue []string

			if reSpace.MatchString(tmpStr) {
				// Split according to space
				words := reSpace.Split(tmpStr, -1)
				// Complete words with splitted space(s)
				i := 0 // To kept the last word index in memory
				for i = 0; i < len(words)-1; i++ {
					// NOTE: All kind of space character will be converted into regular spaces
					subValue = append(subValue, words[i]+" ")
				}
				subValue = append(subValue, words[i])
			} else if reComa.MatchString(tmpStr) {
				// Split according to coma
				words := reComa.Split(tmpStr, -1)
				// Complete words with splitted coma(s)
				i := 0 // To kept the last word index in memory
				for i = 0; i < len(words)-1; i++ {
					// NOTE: All kind of space character will be converted into regular spaces
					subValue = append(subValue, words[i]+",")
				}
				subValue = append(subValue, words[i])
			} else {
				subValue = append(subValue, tmpStr)

			}

			// Now fill the final string, split by length if necessary
			finalStr := prefix
			wi := 0
			lenAvail := lenRem // Available length
			for wi < len(subValue) {
				if lenAvail == 0 {
					finalStr += "\n" + prefix
					lenAvail = lenRem
				}
				// Split the word if longer than the full
				if len(subValue[wi]) > lenRem {
					// Terminate the current line with a fragment of the word
					finalStr += subValue[wi][0:(lenAvail)] + "\n" + prefix
					subValue[wi] = subValue[wi][(lenAvail):]
					lenAvail = lenRem
				} else {
					// The current word is small enough
					// Check if it can be added
					if len(subValue[wi]) <= lenAvail {
						finalStr += subValue[wi]
						lenAvail -= len(subValue[wi])
						wi++
					} else {
						// Remaining place is not enough go to next line
						finalStr += "\n" + prefix
						lenAvail = lenRem
					}
				}
			}

			return finalStr
		}
	}
}
*/
