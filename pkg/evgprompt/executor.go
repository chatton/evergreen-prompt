package evgprompt

import (
	"chatton.com/evergreen-prompt/pkg/evergreen/client"
	"chatton.com/evergreen-prompt/pkg/evergreen/patch"
	"chatton.com/evergreen-prompt/pkg/util/flagutil"
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

func (e *Executor) Execute(in string) {
	if in == "" {
		return
	}
	if in == "quit" || in == "exit" {
		fmt.Printf("Bye!")
		os.Exit(0)
	}

	//if strings.HasPrefix(in, "set-project") {
	//	split := strings.Split(in, " ")
	//	project := split[1]
	//	fmt.Println("setting active project to: " + project)
	//	e.client.ActiveProject = project
	//	return
	//}

	if strings.HasPrefix(in, "patch") {

		args := []string{
			"patch", "-f", "-y",
		}

		if project := flagutil.GetProjectValue(in); project != "" {
			args = append(args, "-p", project)
		}

		if flagutil.HasSpecifiedUncommitted(in) {
			args = append(args, "-u")
		}

		task := flagutil.GetTaskValue(in)
		if task == "" {
			fmt.Println("Task must be specified!")
		}

		args = append(args, "-t", task)

		buildvariant := flagutil.GetBuildVariantValue(in)
		if buildvariant == "" {
			fmt.Println("Buildvariant must be specified!")
		}

		args = append(args, "-v", buildvariant)

		description := flagutil.GetDescriptionValue(in)
		if description == "" {
			description = "evergreen-prompt task"
		}

		args = append(args, "-d", description)

		out, err := exec.Command("evergreen", args...).Output()
		if err != nil {
			fmt.Println(string(out))
			panic(err)
		}
		fmt.Println(string(out))

		id := getPatchIdFromCliOutput(string(out))

		if priority := flagutil.GetPriorityValue(in); priority != "" {
			// set priority of patch
			_, err := e.client.PatchPatch(id, patch.Body{Priority: 10})
			if err != nil {
				panic(err)
			}
		}

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
