package lib

import "strings"

func IsColor(str string) bool {
	if strings.Index(str, "[E]") > -1 {
		return true
	}
	return false
}
