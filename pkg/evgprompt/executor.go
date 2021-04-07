package evgprompt

import (
	"chatton.com/evergreen-prompt/pkg/evergreen/client"
	"chatton.com/evergreen-prompt/pkg/evergreen/patch"
	"chatton.com/evergreen-prompt/pkg/util/flags"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Executor struct {
	client *client.EvergreenClient
}

func NewExecutor(client *client.EvergreenClient) *Executor {
	return &Executor{client: client}
}

func (e *Executor) handleEvergreenPatchCreate(s string) {
	args := []string{
		"patch", "-f", "-y",
	}

	if project := flags.GetProjectValue(s); project != "" {
		args = append(args, "-p", project)
	} else if e.client.DefaultProject != "" {
		args = append(args, "-p", e.client.DefaultProject)
	}

	if flags.HasSpecifiedUncommitted(s) {
		args = append(args, "-u")
	}

	task := flags.GetTaskValue(s)
	if task == "" {
		fmt.Println("Task must be specified!")
	}

	args = append(args, "-t", task)

	buildvariant := flags.GetBuildVariantValue(s)
	if buildvariant == "" {
		fmt.Println("Buildvariant must be specified!")
	}

	args = append(args, "-v", buildvariant)

	description := flags.GetDescriptionValue(s)
	if description == "" {
		description = "evergreen-prompt task"
	}

	args = append(args, "-d", description)

	out, err := exec.Command("evergreen", args...).Output()
	fmt.Println(string(out))
	if err != nil {
		return
	}

	if priority := flags.GetPriorityValue(s); priority != "" {
		p, err := strconv.Atoi(priority)
		if err != nil {
			fmt.Printf("could not convert priority [%s] to an integer!\n", priority)
			return
		}

		id := getPatchIdFromCliOutput(string(out))
		// set priority of patch
		_, err = e.client.PatchPatch(id, patch.Body{Priority: p})
		if err != nil {
			fmt.Printf("error updating patch priority: %s\n", err)
			return
		}
	}

}

func (e *Executor) Execute(in string) {
	if in == "" {
		return
	}
	if in == "quit" || in == "exit" {
		fmt.Printf("Bye!\n")
		os.Exit(0)
	}

	if strings.HasPrefix(in, "patch create") {
		e.handleEvergreenPatchCreate(in)
	}
}

func getPatchIdFromCliOutput(output string) string {
	pattern := fmt.Sprintf(`\s+ID\s:\s([a-zA-Z0-9]+)`)
	r := regexp.MustCompile(pattern)
	allMatches := r.FindStringSubmatch(output)
	if len(allMatches) != 2 {
		return ""
	}
	return allMatches[1]
}
