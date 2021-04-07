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
	flags := extractFlags(s, patchCreate)
	if bv, ok := getValueFromFlagKey("--buildvariant", flags); ok {
		return bv
	}
	return ""
}

func GetTaskValue(s string) string {
	flags := extractFlags(s, patchCreate)
	if task, ok := getValueFromFlagKey("--task", flags); ok {
		return task
	}
	return ""
}

func GetDescriptionValue(s string) string {
	flags := extractFlags(s, patchCreate)
	if description, ok := getValueFromFlagKey("--description", flags); ok {
		return description
	}
	return ""
}

func GetPriorityValue(s string) string {
	flags := extractFlags(s, patchCreate)
	if priority, ok := getValueFromFlagKey("--priority", flags); ok {
		return priority
	}
	return ""
}

func GetProjectValue(s string) string {
	flags := extractFlags(s, patchCreate)
	if project, ok := getValueFromFlagKey("--project", flags); ok {
		return project
	}
	return ""
}

func HasSpecifiedUncommitted(s string) bool {
	flags := extractFlags(s, patchCreate)
	_, ok := getValueFromFlagKey("--uncommitted", flags)
	return ok
}

func GetAllParams(s string) []string {
	flags := extractFlags(s, patchCreate)
	var paramFlags []string
	for _, f := range flags {
		if f.key == "--param" {
			paramFlags = append(paramFlags, f.value)
		}
	}
	return paramFlags
}

func getValueFromFlagKey(key string, flags []flag) (string, bool) {
	for _, f := range flags {
		if f.key == key {
			return f.value, true
		}
	}
	return "", false
}

type flag struct {
	key   string
	value string
}

// extractFlags converts a string with the given prefix into a map[string]string
// with the keys as the flags and the provided values as the map values.
// if the flag does not require a value, an empty string will be set as the value
// in the map.
func extractFlags(s, prefix string) []flag {
	var flags []flag

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
		flags = append(flags, flag{
			key:   args[i],
			value: args[i+1],
		})
	}

	return flags
}
