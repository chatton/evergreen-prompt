package flagutil

import (
	"fmt"
	"regexp"
	"strings"
)

func GetBuildVariantValue(s string) string {
	return ExtractFlagValue("--buildvariant", s)
}

func GetTaskValue(s string) string {
	return ExtractFlagValue("--task", s)
}

func GetDescriptionValue(s string) string {
	return ExtractFlagValue("--description", s)
}

func GetPriorityValue(s string) string {
	return ExtractFlagValue("--priority", s)
}

func ExtractFlagValue(flag, text string) string {
	pattern := fmt.Sprintf(`%s\s+(["a-z0-9_\s]+)`, flag)
	r := regexp.MustCompile(pattern)
	allMatches := r.FindStringSubmatch(text)

	if len(allMatches) != 2 {
		return ""
	}
	return strings.TrimRight(allMatches[1], " ")
}
