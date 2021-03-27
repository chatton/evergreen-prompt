package evergreen

import (
	"github.com/goccy/go-yaml"
	"io/ioutil"
	"path/filepath"
)

type Configuration struct {
	BuildVariants []BuildVariant `yaml:"buildvariants"`
	Tasks         []Task         `yaml:"tasks"`
	TaskGroups    []TaskGroup    `yaml:"task_groups"`
}

type Task struct {
	Name string `yaml:"name"`
}

type TaskGroup struct {
	Name  string   `yaml:"name"`
	Tasks []string `yaml:"tasks"`
}

type BuildVariant struct {
	// Each task can be the name of a task or a task group
	Tasks []Task `yaml:"tasks"`
	Name  string `yaml:"name"`
}

// GetTasksInBuildVariant returns all of the Tasks that are associated with
// the given buildvariant.
func (c Configuration) GetTasksInBuildVariant(buildVariantName string) []Task {
	var tasks []Task
	for _, bv := range c.BuildVariants {
		if bv.Name != buildVariantName {
			continue
		}
		for _, task := range bv.Tasks {
			// the task in the buildvariant is a task group, not a task.
			if tg := c.getTaskGroupByName(task.Name); tg != nil {
				for _, taskName := range tg.Tasks {
					tasks = append(tasks, Task{Name: taskName})
				}
			} else {
				tasks = append(tasks, task)
			}
		}
	}

	return tasks
}

func (c Configuration) getTaskGroupByName(taskName string) *TaskGroup {
	for _, tg := range c.TaskGroups {
		if tg.Name == taskName {
			return &tg
		}
	}
	return nil
}

// GetBuildVariantsThatTaskIsIn returns a list of buildvariants that reference
// the given task.
func (c Configuration) GetBuildVariantsThatTaskIsIn(taskName string) []BuildVariant {
	var buildVariants []BuildVariant
	for _, bv := range c.BuildVariants {
		bvTasks := c.GetTasksInBuildVariant(bv.Name)
		if containsTask(bvTasks, taskName) {
			buildVariants = append(buildVariants, bv)
		}
	}
	return buildVariants
}

func containsTask(tasks []Task, taskName string) bool {
	for _, t := range tasks {
		if t.Name == taskName {
			return true
		}
	}
	return false
}

func FromYamlConfigFile() (Configuration, error) {
	evgFilePath, err := filepath.Abs(".evergreen.yml")
	if err != nil {
		return Configuration{}, err
	}
	config := Configuration{}
	bytes, err := ioutil.ReadFile(evgFilePath)
	if err != nil {
		return Configuration{}, err
	}

	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return Configuration{}, err
	}

	return config, nil
}
