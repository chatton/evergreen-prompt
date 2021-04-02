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

func GetProjectValue(s string) string {
	return ExtractFlagValue("--project", s)
}

func HasSpecifiedUncommitted(s string) bool {
	return strings.Contains(s, "--uncommitted")
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

// ExtractFlags converts a string with the given prefix into a map[string]string
// with the keys as the flags and the provided values as the map values.
// if the flag does not require a value, an empty string will be set as the value
// in the map/
func ExtractFlags(s, prefix string) map[string]string {
	flags := map[string]string{}

	flagsOnly := strings.TrimLeft(s, prefix)
	flagsOnly = strings.TrimSpace(flagsOnly)
	splitString := strings.Split(flagsOnly, " ")

	var args []string
	prevElement := ""
	for _, s := range splitString {

		// we see two flags in a row, which means we need to insert an empty string for the previous
		// flag.
		if strings.HasPrefix(prevElement, "--") && strings.HasPrefix(s, "--") {
			args = append(args, "", s)
			prevElement = s
			continue
		}

		args = append(args, s)
		prevElement = s
	}

	// if we have an odd number of elements, it means the last item was a flag that does
	// not accept an argument, let's add an empty string as its argument.
	if len(args)%2 != 0 {
		args = append(args, "")
	}

	// convert every element into a map
	for i := 0; i < len(args); i += 2 {
		flags[args[i]] = args[i+1]
	}

	return flags
}
