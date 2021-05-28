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

func (c *Completer) createPatchSuggestions(d prompt.Document, inputFlags flags.Flags) []prompt.Suggest {

	if getLastWord(d) == "--project" {
		return c.getProjectSuggestions(d)
	}

	if getLastWord(d) == "--task" {
		return c.getTaskSuggestions(d, inputFlags)
	}

	if getLastWord(d) == "--buildvariant" {
		return c.getBuildVariantSuggestions(d, inputFlags)
	}

	if getLastWord(d) == "--param" {
		return c.getParamSuggestions()
	}

	if getLastWord(d) == "--description" {
		return nil
	}

	if getLastWord(d) == "--priority" {
		return nil
	}

	var suggestions []prompt.Suggest

	suggestions = append(suggestions,
		prompt.Suggest{
			Text:        "--task",
			Description: "Specify a task to run",
		})

	suggestions = append(suggestions,
		prompt.Suggest{
			Text:        "--buildvariant",
			Description: "Specify a build variant",
		})

	if inputFlags.Description == "" {
		suggestions = append(suggestions,
			prompt.Suggest{
				Text:        "--description",
				Description: "Specify a description for the patch",
			})
	}

	if inputFlags.Priority == -1 {
		suggestions = append(suggestions,
			prompt.Suggest{
				Text:        "--priority",
				Description: "Specify the priority for the patch",
			})
	}

	if !inputFlags.Uncommitted {
		suggestions = append(suggestions,
			prompt.Suggest{
				Text:        "--uncommitted",
				Description: "Include uncommitted changes",
			})
	}

	if inputFlags.Project == "" {
		suggestions = append(suggestions,
			prompt.Suggest{
				Text:        "--project",
				Description: "Specify the name of an existing evergreen project",
			})
	}

	// it's possible to add multiple params to a single command, so the field
	// already existing doesn't mean we shouldn't show it.
	suggestions = append(suggestions, prompt.Suggest{
		Text:        "--param",
		Description: "Specify a parameter for the evergreen patch",
	})

	return prompt.FilterFuzzy(suggestions, d.GetWordBeforeCursor(), true)

}

func (c *Completer) patchSuggestions(d prompt.Document) []prompt.Suggest {

	text := d.TextBeforeCursor()
	inputFlags, err := flags.Parse(text)
	if err != nil {
		panic(err)
	}
	if strings.Contains(text, "create") {
		return c.createPatchSuggestions(d, inputFlags)
	}

	return prompt.FilterFuzzy([]prompt.Suggest{
		{
			Text:        "create",
			Description: "Create a new patch",
		},
	}, d.GetWordBeforeCursor(), true)
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

func (c *Completer) getTaskSuggestions(d prompt.Document, flags flags.Flags) []prompt.Suggest {
	var suggestions []prompt.Suggest

	// if we are getting the task and the buildvariant already specified, we need to show
	// only the tasks that contain this build varient otherwise we can show all the tasks.
	buildvariantValue := ""
	if len(flags.BuildVariants) > 0 {
		buildvariantValue = flags.BuildVariants[len(flags.BuildVariants)-1]
	}

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

func (c *Completer) getBuildVariantSuggestions(d prompt.Document, flags flags.Flags) []prompt.Suggest {
	var suggestions []prompt.Suggest

	// if we are getting the build variant and the task is already specified, we need to show
	// only the build variants that contain this task, otherwise we can show all the buildvariants.
	taskValue := ""
	if len(flags.Tasks) > 0 {
		taskValue = flags.Tasks[len(flags.Tasks)-1]
	}

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

func (c *Completer) getParamSuggestions() []prompt.Suggest {
	var suggestions []prompt.Suggest

	for _, param := range c.config.Parameters {
		suggestions = append(suggestions, prompt.Suggest{
			Text:        param.Key + "=" + param.Value,
			Description: param.Description,
		})
	}

	return suggestions
}
