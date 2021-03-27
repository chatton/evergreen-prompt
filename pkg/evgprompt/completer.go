package evgprompt

import (
	"chatton.com/evergreen-prompt/pkg/evergreen/client"
	"github.com/c-bata/go-prompt"
	"strings"
)

func NewCompleter(c *client.EvergreenClient) Completer {
	return Completer{
		client: c,
	}
}

type Completer struct {
	client *client.EvergreenClient
}

func (c *Completer) Complete(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	if suggestions, found := c.completeOptionArguments(d); found {
		return suggestions
	}

	// don't display set-project if it has already been set.
	if strings.Contains(d.TextBeforeCursor(), "set-project") {
		return nil
	}

	return prompt.FilterFuzzy([]prompt.Suggest{
		{
			Text:        "set-project",
			Description: "Choose the active project",
		},
	}, d.GetWordBeforeCursor(), true,
	)
}


func getPreviousOption(d prompt.Document) (cmd, option string, found bool) {
	args := strings.Split(d.TextBeforeCursor(), " ")
	l := len(args)
	if l >= 2 {
		option = args[l-2]
		return "", option, true
	}
	if strings.HasPrefix(option, "-") {
		return args[0], option, true
	}
	return "", "", false
}

func (c *Completer) completeOptionArguments(d prompt.Document) ([]prompt.Suggest, bool) {
	_, previousWord, found := getPreviousOption(d)
	if !found {
		return []prompt.Suggest{}, false
	}

	if previousWord == "set-project" {
		return prompt.FilterFuzzy(c.projectSuggestions(), d.GetWordBeforeCursor(), true), true
	}

	return []prompt.Suggest{}, false
}

func (c *Completer) projectSuggestions() []prompt.Suggest {
	projects, err := c.client.GetProjects()
	if err != nil {
		panic(err)
	}

	var suggestions []prompt.Suggest
	for _, p := range projects {
		suggestions = append(suggestions, prompt.Suggest{
			Text: p,
		})
	}

	return suggestions
}
