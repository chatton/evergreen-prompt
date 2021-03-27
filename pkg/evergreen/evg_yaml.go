package evergreen

import (
	"github.com/goccy/go-yaml"
	"io/ioutil"
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
	Tasks []string `yaml:"tasks"`
}

type BuildVariant struct {
	// Each task can be the name of a task or a task group
	Tasks []Task `yaml:"tasks"`
}

func FromYamlConfigFile(fullPathOfEvergreenYaml string) (Configuration, error) {
	config := Configuration{}
	bytes, err := ioutil.ReadFile(fullPathOfEvergreenYaml)
	if err != nil {
		return Configuration{}, err
	}

	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return Configuration{}, err
	}

	return config, nil
}
