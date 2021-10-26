package feature

import (
	"fmt"
	"regexp"
)

const (
	D_RAWSTR    string = "NO_VALUE"
	D_ISBOOL    bool   = false
	D_HASQUOTE  bool   = true
	D_TAGPREFIX string = "/"
)

type Value struct {
	RawStr   string
	IsBool   bool
	HasQuote bool
}

func NewValue(s ...string) *Value {
	if s == nil {
		return &Value{D_RAWSTR, D_ISBOOL, D_HASQUOTE}
	} else {
		if len(s) > 1 {
			panic("You are supposed to provide only one string value.")
		}
		return &Value{s[0], D_ISBOOL, D_HASQUOTE}
	}
}

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
