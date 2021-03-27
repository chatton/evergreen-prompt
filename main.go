package main

import (
	"chatton.com/evergreen-prompt/pkg/evergreen"
	"chatton.com/evergreen-prompt/pkg/evergreen/client"
	"chatton.com/evergreen-prompt/pkg/evgprompt"
	"fmt"
	"github.com/c-bata/go-prompt"
)

func main() {

	config, err := evergreen.FromYamlConfigFile("/Users/cianhatton/go/src/github.com/mongodb/mongodb-kubernetes-operator/.evergreen.yml")
	if err != nil {
		panic(err)
	}
	defer fmt.Println("Bye!")

	c, err := client.NewEvergreenClient()
	if err != nil {
		panic(err)
	}

	completer := evgprompt.NewCompleter(c, config)
	executor := evgprompt.NewExecutor(c)

	p := prompt.New(executor.Execute, completer.Complete,
		prompt.OptionTitle("evergreen-prompt"))
		prompt.OptionPrefix(">>> ")
	p.Run()
}
