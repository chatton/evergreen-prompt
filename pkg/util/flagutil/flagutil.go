package flagutil

import (
	"fmt"
	"regexp"
)

func GetBuildVariantFlag(s string) string {
	return ExtractFlagValue("--buildvariant", s)
}

func GetTaskFlag(s string) string {
	return ExtractFlagValue("--task", s)
}

func ExtractFlagValue(flag, text string) string {
	pattern := fmt.Sprintf(`%s\s([a-z0-9_]+)`, flag)
	r := regexp.MustCompile(pattern)
	allMatches := r.FindStringSubmatch(text)
	if len(allMatches) != 2 {
		return ""
	}

	return allMatches[1]
}
