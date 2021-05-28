package evgprompt

import (
	"chatton.com/evergreen-prompt/pkg/evergreen/client"
	"chatton.com/evergreen-prompt/pkg/evergreen/patch"
	"chatton.com/evergreen-prompt/pkg/util/flags"
	"fmt"
	"os"
	"os/exec"
	"regexp"
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

	input, err := flags.Parse(s)
	if err != nil {
		fmt.Println(err)
		return
	}

	if input.Project != "" {
		args = append(args, "-p", input.Project)
	} else if e.client.DefaultProject != "" {
		args = append(args, "-p", e.client.DefaultProject)
	}

	if input.Uncommitted {
		args = append(args, "-u")
	}

	if input.Tasks == nil {
		fmt.Println("Task mut be specified!")
		return
	}

	// specify each individual task as a separate argument.
	for _, t := range input.Tasks {
		args = append(args, "-t", t)
	}

	// p is expected to be in the form of "Key=Value"
	for k, v := range input.Params {
		args = append(args, "--param", k+"="+v)
	}

	for _, bv := range input.BuildVariants {
		args = append(args, "-v", bv)
	}

	if len(input.BuildVariants) == 0 {
		fmt.Println("Buildvariant must be specified!")
		return
	}

	description := input.Description
	if input.Description == "" {
		description = "evergreen-prompt task"
	}

	args = append(args, "-d", description)

	out, err := exec.Command("evergreen", args...).Output()
	fmt.Println(string(out))
	if err != nil {
		return
	}

	// -1 means we didn't set a priority. Just use default.
	if input.Priority != -1 {
		id := getPatchIdFromCliOutput(string(out))
		// set priority of patch
		_, err = e.client.PatchPatch(id, patch.Body{Priority: input.Priority})
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
