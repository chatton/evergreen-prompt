package flagutil

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"regexp"
)

func GetBuildVariantFlag(d prompt.Document) string {
	return ExtractFlagValue("--buildvariant", d.TextBeforeCursor())
}

func GetTaskFlag(d prompt.Document) string {
	return ExtractFlagValue("--task", d.TextBeforeCursor())
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
