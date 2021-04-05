package evgprompt

import (
	"chatton.com/evergreen-prompt/pkg/evergreen"
	"chatton.com/evergreen-prompt/pkg/evergreen/client"
	"chatton.com/evergreen-prompt/pkg/util/flags"
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

// isPatchCommand returns true if the current command is a
// patch command.
func isPatchCommand(d prompt.Document) bool {
	text := d.TextBeforeCursor()
	text = strings.TrimSpace(text)
	return strings.HasPrefix(text, "patch")
}

func (c *Completer) Complete(d prompt.Document) []prompt.Suggest {
	if isPatchCommand(d) {
		return c.patchSuggestions(d)
	}

	return prompt.FilterFuzzy([]prompt.Suggest{
		{
			Text:        "patch",
			Description: "run an evergreen patch",
		},
	}, d.GetWordBeforeCursor(), true,
	)
}

func (c *Completer) patchSuggestions(d prompt.Document) []prompt.Suggest {

	if getLastWord(d) == "--project" {
		return c.getProjectSuggestions(d)
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

	var suggestions []prompt.Suggest

	// we only want to show suggestions when they have not yet been specified.
	if flags.GetTaskValue(d.TextBeforeCursor()) == "" {
		suggestions = append(suggestions,
			prompt.Suggest{
				Text:        "--task",
				Description: "Specify a task to run",
			})
	}

	if flags.GetBuildVariantValue(d.TextBeforeCursor()) == "" {
		suggestions = append(suggestions,
			prompt.Suggest{
				Text:        "--buildvariant",
				Description: "Specify a build variant",
			})
	}

	if flags.GetDescriptionValue(d.TextBeforeCursor()) == "" {
		suggestions = append(suggestions,
			prompt.Suggest{
				Text:        "--description",
				Description: "Specify a description for the patch",
			})
	}

	if flags.GetPriorityValue(d.TextBeforeCursor()) == "" {
		suggestions = append(suggestions,
			prompt.Suggest{
				Text:        "--priority",
				Description: "Specify the priority for the patch",
			})
	}

	if !flags.HasSpecifiedUncommitted(d.TextBeforeCursor()) {
		suggestions = append(suggestions,
			prompt.Suggest{
				Text:        "--uncommitted",
				Description: "Include uncommitted changes",
			})
	}

	if flags.GetProjectValue(d.TextBeforeCursor()) == "" {
		suggestions = append(suggestions,
			prompt.Suggest{
				Text:        "--project",
				Description: "Specify the name of an existing evergreen project",
			})
	}

	return prompt.FilterFuzzy(suggestions, d.GetWordBeforeCursor(), true)
}

func (c *Completer) getProjectSuggestions(d prompt.Document) []prompt.Suggest {
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

	return prompt.FilterFuzzy(suggestions, d.GetWordBeforeCursor(), true)

}

func (c *Completer) getTaskSuggestions(d prompt.Document) []prompt.Suggest {
	var suggestions []prompt.Suggest

	// if we are getting the task and the buildvariant already specified, we need to show
	// only the tasks that contain this build varient otherwise we can show all the tasks.
	buildvariantValue := flags.GetBuildVariantValue(d.TextBeforeCursor())

	if buildvariantValue == "" {
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

	return prompt.FilterFuzzy(suggestions, d.GetWordBeforeCursor(), true)
}

func (c *Completer) getBuildVariantSuggestions(d prompt.Document) []prompt.Suggest {
	var suggestions []prompt.Suggest

	// if we are getting the build variant and the task is already specified, we need to show
	// only the build variants that contain this task, otherwise we can show all the buildvariants.
	taskValue := flags.GetTaskValue(d.TextBeforeCursor())

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

	return prompt.FilterFuzzy(suggestions, d.GetWordBeforeCursor(), true)
}
