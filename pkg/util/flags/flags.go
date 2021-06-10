package flags

import (
	"regexp"
	"strconv"
	"strings"
)

// pattern used to split flags from input string.
var flagsPattern *regexp.Regexp

const patchCreate = "patch create"

func init() {
	// split on spaces but not when between quotes, this allows for the description
	// field or any quoted fields to treated the same way.
	// https://stackoverflow.com/questions/47489745/splitting-a-string-at-space-except-inside-quotation-marks/47490643#47490643
	flagsPattern = regexp.MustCompile(`[^\s"]+|"([^"]*)"`)
}

// Flags stores all of the values parsed from the command line flags
type Flags struct {
	BuildVariants []string
	Tasks         []string
	Description   string
	Priority      int
	Project       string
	Uncommitted   bool
	Params        map[string]string
	Times         int
}

// Parse accepts the command line input string and returns a struct
// containing all of the values that were specified.
func Parse(in string) (Flags, error) {
	priorityStr := getPriorityValue(in)
	if priorityStr == "" {
		priorityStr = "-1"
	}

	priority, err := strconv.Atoi(priorityStr)
	if err != nil {
		return Flags{}, err
	}

	return Flags{
		BuildVariants: getAllBuildVariants(in),
		Tasks:         getAllTasks(in),
		Description:   getDescriptionValue(in),
		Project:       getProjectValue(in),
		Priority:      priority,
		Uncommitted:   hasSpecifiedUncommitted(in),
		Params:        getAllParams(in),
		Times:         getTimes(in),
	}, nil
}

func getTimes(s string) int {
	flags := extractFlags(s, patchCreate)
	if count, ok := getValuesWithFlagKey("--times", flags); ok {
		times, err := strconv.Atoi(count[0])
		if err != nil {
			return -1
		}

		return times
	}
	return -1
}

func getBuildVariantValue(s string) string {
	flags := extractFlags(s, patchCreate)
	if bv, ok := getValueFromFlagKey("--buildvariant", flags); ok {
		return bv
	}
	return ""
}

func getAllBuildVariants(s string) []string {
	flags := extractFlags(s, patchCreate)
	if bvs, ok := getValuesWithFlagKey("--buildvariant", flags); ok {
		return bvs
	}
	return nil
}

func getTaskValue(s string) string {
	flags := extractFlags(s, patchCreate)
	if task, ok := getValueFromFlagKey("--task", flags); ok {
		return task
	}
	return ""
}

func getAllTasks(s string) []string {
	flags := extractFlags(s, patchCreate)
	if tasks, ok := getValuesWithFlagKey("--task", flags); ok {
		return tasks
	}
	return nil
}

func getDescriptionValue(s string) string {
	flags := extractFlags(s, patchCreate)
	if description, ok := getValueFromFlagKey("--description", flags); ok {
		return description
	}
	return ""
}

func getPriorityValue(s string) string {
	flags := extractFlags(s, patchCreate)
	if priority, ok := getValueFromFlagKey("--priority", flags); ok {
		return priority
	}
	return ""
}

func getProjectValue(s string) string {
	flags := extractFlags(s, patchCreate)
	if project, ok := getValueFromFlagKey("--project", flags); ok {
		return project
	}
	return ""
}

func hasSpecifiedUncommitted(s string) bool {
	flags := extractFlags(s, patchCreate)
	_, ok := getValueFromFlagKey("--uncommitted", flags)
	return ok
}

func getAllParams(s string) map[string]string {
	flags := extractFlags(s, patchCreate)
	params := map[string]string{}
	for _, f := range flags {
		if f.key == "--param" {
			params[f.key] = f.value
		}
	}
	return params
}

func getValueFromFlagKey(key string, flags []flag) (string, bool) {
	for _, f := range flags {
		if f.key == key {
			return f.value, true
		}
	}
	return "", false
}

func getValuesWithFlagKey(key string, flags []flag) ([]string, bool) {
	var results []string
	for _, f := range flags {
		if f.key == key {
			results = append(results, f.value)
		}
	}
	return results, len(results) > 0
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
