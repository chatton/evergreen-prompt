package evgprompt

import (
	"chatton.com/evergreen-prompt/pkg/evergreen"
	"chatton.com/evergreen-prompt/pkg/evergreen/client"
	"chatton.com/evergreen-prompt/pkg/util/flagutil"
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
	if getLastWord(d) == "--project" {
		return prompt.FilterFuzzy(c.projectSuggestions(), d.GetWordBeforeCursor(), true)
	}

	if getLastWord(d) == "--task" {
		return c.getTaskSuggestions(d)
	}

	if getLastWord(d) == "--buildvariant" {
		return c.getBuildVariantSuggestions(d)
	}

	if getLastWord(d) == "--description" {
		return nil
	}

	if getLastWord(d) == "--priority" {
		return nil
	}

	//if getLastWord(d) == "--uncommitted" {
	//	return nil
	//}

	if strings.HasPrefix(d.TextBeforeCursor(), "patch") {
		return patchSuggestions(d)
	}

	// don't display set-project if it has already been set.
	if strings.Contains(d.TextBeforeCursor(), "set-project") {
		return nil
	}

	return prompt.FilterFuzzy([]prompt.Suggest{
		{
			Text:        "patch",
			Description: "run an evergreen patch",
		},
	}, d.GetWordBeforeCursor(), true,
	)
}

func patchSuggestions(d prompt.Document) []prompt.Suggest {
	var suggestions []prompt.Suggest

	// we only want to show suggestions when they have not yet been specified.
	if flagutil.GetTaskValue(d.TextBeforeCursor()) == "" {
		suggestions = append(suggestions,
			prompt.Suggest{
				Text:        "--task",
				Description: "Specify a task to run",
			})
	}

	if flagutil.GetBuildVariantValue(d.TextBeforeCursor()) == "" {
		suggestions = append(suggestions,
			prompt.Suggest{
				Text:        "--buildvariant",
				Description: "Specify a build variant",
			})
	}

	if flagutil.GetDescriptionValue(d.TextBeforeCursor()) == "" {
		if !strings.Contains(d.TextBeforeCursor(), "--description") {
			suggestions = append(suggestions,
				prompt.Suggest{
					Text:        "--description",
					Description: "Specify a description for the patch",
				})
		}
	}

	if flagutil.GetPriorityValue(d.TextBeforeCursor()) == "" {
		if !strings.Contains(d.TextBeforeCursor(), "--priority") {
			suggestions = append(suggestions,
				prompt.Suggest{
					Text:        "--priority",
					Description: "Specify the priority for the patch",
				})
		}
	}

	if !flagutil.HasSpecifiedUncommitted(d.TextBeforeCursor()) {
		suggestions = append(suggestions,
			prompt.Suggest{
				Text:        "--uncommitted",
				Description: "Include uncommitted changes",
			})
	}

	if flagutil.GetProjectValue(d.TextBeforeCursor()) == "" {
		suggestions = append(suggestions,
			prompt.Suggest{
				Text:        "--project",
				Description: "Specify the name of an existing evergreen project",
			})
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

func (c *Completer) getTaskSuggestions(d prompt.Document) []prompt.Suggest {
	var suggestions []prompt.Suggest

	// if we are getting the task and the buildvariant already specified, we need to show
	// only the tasks that contain this build varient otherwise we can show all the tasks.
	buildvariantValue := flagutil.GetBuildVariantValue(d.TextBeforeCursor())

	if buildvariantValue == ""  {
		for _, t := range c.config.Tasks {
			suggestions = append(suggestions, prompt.Suggest{
				Text: t.Name,
			})
		}
		return prompt.FilterFuzzy(suggestions, d.GetWordBeforeCursor(), true)
	}

	for _, t := range c.config.GetTasksInBuildVariant(buildvariantValue) {
		suggestions = append(suggestions, prompt.Suggest{
			Text: t.Name,
		})
	}

	//suggestions = append(suggestions, prompt.Suggest{
	//	Text: "all",
	//})

	return prompt.FilterFuzzy(suggestions, d.GetWordBeforeCursor(), true)
}

func (c *Completer) getBuildVariantSuggestions(d prompt.Document) []prompt.Suggest {
	var suggestions []prompt.Suggest

	// if we are getting the build variant and the task is already specified, we need to show
	// only the build variants that contain this task, otherwise we can show all the buildvariants.
	taskValue := flagutil.GetTaskValue(d.TextBeforeCursor())

	if taskValue == "" {
		for _, bv := range c.config.BuildVariants {
			suggestions = append(suggestions, prompt.Suggest{
				Text: bv.Name,
			})
		}
		return prompt.FilterFuzzy(suggestions, d.GetWordBeforeCursor(), true)
	}

	for _, bv := range c.config.GetBuildVariantsThatTaskIsIn(taskValue) {
		suggestions = append(suggestions, prompt.Suggest{
			Text: bv.Name,
		})
	}

	//suggestions = append(suggestions, prompt.Suggest{
	//	Text: "all",
	//})

	return prompt.FilterFuzzy(suggestions, d.GetWordBeforeCursor(), true)
}
