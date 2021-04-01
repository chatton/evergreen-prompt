package evgprompt

import (
	"chatton.com/evergreen-prompt/pkg/evergreen/client"
	"chatton.com/evergreen-prompt/pkg/util/flagutil"
	"fmt"
	"os"
	"os/exec"
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

	if strings.HasPrefix(in, "set-project") {
		split := strings.Split(in, " ")
		project := split[1]
		fmt.Println("setting active project to: " + project)
		e.client.ActiveProject = project
		return
	}

	if strings.HasPrefix(in, "patch") {

		task := flagutil.GetTaskValue(in)
		if task == "" {
			fmt.Println("Task must be specified!")
		}

		buildvariant := flagutil.GetBuildVariantValue(in)
		if buildvariant == "" {
			fmt.Println("Buildvariant must be specified!")
		}

		description := flagutil.GetDescriptionValue(in)
		if description == "" {
			description = "evergreen-prompt task"
		}

		out, err := exec.Command("evergreen", "patch", "-p", e.client.ActiveProject, "-f", "-u", "-d", description, "-t", task, "-v", buildvariant, "-y").Output()
		if err != nil {
			panic(err)
		}
		fmt.Println(string(out))
		//cmd.Stdin = os.Stdin
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		//if err := cmd.Run(); err != nil {
		//	panic(err)
		//}

		//return
		//}
	}
}
