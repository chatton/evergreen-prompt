package flagutil

import (
	"regexp"
	"strings"
)

// pattern used to split flags from input string.
var flagsPattern *regexp.Regexp

func init() {
	// split on spaces but not when between quotes, this allows for the description
	// field or any quoted fields to treated the same way.
	// https://stackoverflow.com/questions/47489745/splitting-a-string-at-space-except-inside-quotation-marks/47490643#47490643
	flagsPattern = regexp.MustCompile(`[^\s"]+|"([^"]*)"`)
}

func GetBuildVariantValue(s string) string {
	if bv, ok := ExtractFlags(s, "patch")["--buildvariant"]; ok {
		return bv
	}
	return ""
}

func GetTaskValue(s string) string {
	if task, ok := ExtractFlags(s, "patch")["--task"]; ok {
		return task
	}
	return ""
}

func GetDescriptionValue(s string) string {
	if description, ok := ExtractFlags(s, "patch")["--description"]; ok {
		return description
	}
	return ""
}

func GetPriorityValue(s string) string {
	if priority, ok := ExtractFlags(s, "patch")["--priority"]; ok {
		return priority
	}
	return ""
}

func GetProjectValue(s string) string {
	if project, ok := ExtractFlags(s, "patch")["--project"]; ok {
		return project
	}
	return ""
}

func HasSpecifiedUncommitted(s string) bool {
	_, ok := ExtractFlags(s, "patch")["--uncommitted"]
	return ok
}

// ExtractFlags converts a string with the given prefix into a map[string]string
// with the keys as the flags and the provided values as the map values.
// if the flag does not require a value, an empty string will be set as the value
// in the map/
func ExtractFlags(s, prefix string) map[string]string {
	flags := map[string]string{}

	flagsOnly := strings.TrimLeft(s, prefix)
	flagsOnly = strings.TrimSpace(flagsOnly)
	splitString := flagsPattern.FindAllString(flagsOnly, -1)

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
