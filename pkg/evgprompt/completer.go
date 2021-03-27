package evgprompt

import (
	"chatton.com/evergreen-prompt/pkg/evergreen"
	"chatton.com/evergreen-prompt/pkg/evergreen/client"
	"github.com/c-bata/go-prompt"
	"strings"
)

func NewCompleter(c *client.EvergreenClient, config evergreen.Configuration) Completer {
	return Completer{
		client: c,
		config: config,
	}
}

type Completer struct {
	client *client.EvergreenClient
	config evergreen.Configuration
}

// getLastWord returns rightmost word.
func getLastWord(d prompt.Document) string {
	text := d.TextBeforeCursor()
	args := strings.Split(text, " ")
	if len(args) > 1 {
		return args[len(args)-2]
	}
	return ""
}

func (c *Completer) Complete(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		suggestions := []prompt.Suggest{
			{
				Text:        "patch",
				Description: "",
			},
			{
				Text:        "set-project",
				Description: "",
			},
		}
		return prompt.FilterFuzzy(suggestions, d.GetWordBeforeCursor(), true)
	}

	if getLastWord(d) == "set-project" {
		return prompt.FilterFuzzy(c.projectSuggestions(), d.GetWordBeforeCursor(), true)
	}

	if getLastWord(d) == "--task" {
		suggestions := []prompt.Suggest{
			{
				Text:        "task-1",
				Description: "",
			},
			{
				Text:        "task-2",
				Description: "",
			},
		}
		return prompt.FilterFuzzy(suggestions, d.GetWordBeforeCursor(), true)
	}

	if getLastWord(d) == "--buildvariant" {
		suggestions := []prompt.Suggest{
			{
				Text:        "buildvariant-1",
				Description: "",
			},
			{
				Text:        "buildvariant-2",
				Description: "",
			},
		}
		return prompt.FilterFuzzy(suggestions, d.GetWordBeforeCursor(), true)
	}

	if strings.HasPrefix(d.TextBeforeCursor(), "patch") {
		return patchSuggestions(d)
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

func patchSuggestions(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{
		{
			Text:        "--task",
			Description: "Specify a task to run",
		},
		{
			Text:        "--buildvariant",
			Description: "Specify a build variant",
		},
	}
	return prompt.FilterFuzzy(suggestions, d.GetWordBeforeCursor(), true)
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
