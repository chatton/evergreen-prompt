package evgprompt

import (
	"chatton.com/evergreen-prompt/pkg/evergreen/client"
	"fmt"
	"os"
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
	}
}


