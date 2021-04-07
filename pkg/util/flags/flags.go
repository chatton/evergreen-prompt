package flags

import (
	"regexp"
	"strings"
)

// pattern used to split flags from input string.
var flagsPattern *regexp.Regexp

const patchCreate = "patch create"
const patchAbort = "patch abort"

func init() {
	// split on spaces but not when between quotes, this allows for the description
	// field or any quoted fields to treated the same way.
	// https://stackoverflow.com/questions/47489745/splitting-a-string-at-space-except-inside-quotation-marks/47490643#47490643
	flagsPattern = regexp.MustCompile(`[^\s"]+|"([^"]*)"`)
}

func GetBuildVariantValue(s string) string {
	if bv, ok := extractFlags(s, patchCreate)["--buildvariant"]; ok {
		return bv
	}
	return ""
}

func GetTaskValue(s string) string {
	if task, ok := extractFlags(s, patchCreate)["--task"]; ok {
		return task
	}
	return ""
}

func GetDescriptionValue(s string) string {
	if description, ok := extractFlags(s, patchCreate)["--description"]; ok {
		return description
	}
	return ""
}

func GetPriorityValue(s string) string {
	if priority, ok := extractFlags(s, patchCreate)["--priority"]; ok {
		return priority
	}
	return ""
}

func GetProjectValue(s string) string {
	if project, ok := extractFlags(s, patchCreate)["--project"]; ok {
		return project
	}
	return ""
}

func HasSpecifiedUncommitted(s string) bool {
	_, ok := extractFlags(s, patchCreate)["--uncommitted"]
	return ok
}

func GetPatchId(s string) string {
	if patchId, ok := extractFlags(s, patchAbort)["--patch-id"]; ok {
		return patchId
	}
	return ""
}

// extractFlags converts a string with the given prefix into a map[string]string
// with the keys as the flags and the provided values as the map values.
// if the flag does not require a value, an empty string will be set as the value
// in the map.
func extractFlags(s, prefix string) map[string]string {
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
