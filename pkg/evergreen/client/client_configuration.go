package client

import (
	"encoding/json"
	"fmt"
	"github.com/goccy/go-yaml"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	User           string `json:"User"`
	BaseUrl        string `json:"BaseUrl"`
	ApiKey         string `json:"ApiKey"`
	DefaultProject string `json:"defaultProject"`
}

// LoadConfig loads the config file that contains the User, baseurl and apikey
// required to interact with the evergreen api.
func LoadConfig() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	evergreenConfigPath := path.Join(home, ".evergreen.yml")
	fileBytes, err := ioutil.ReadFile(evergreenConfigPath)

	if err == nil {
		fmt.Println(fmt.Sprintf("Found an existing evergreen configuration in %s", evergreenConfigPath))
		return loadConfigFromEvergreenYaml(fileBytes)
	}

	configFilePath := path.Join(home, ".evergreen-prompt.json")
	bytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("could not read evergreen-prompt.json config file from: %s", configFilePath)
	}

	config := Config{}
	if err := json.Unmarshal(bytes, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

type EvergreenYamlConfiguration struct {
	BaseUrl              string            `yaml:"ui_server_host"`
	ApiKey               string            `yaml:"api_key"`
	User                 string            `yaml:"user"`
	Projects             []Project         `yaml:"projects"`
	ProjectsForDirectory map[string]string `yaml:"projects_for_directory"`
}

type Project struct {
	Name                 string            `yaml:"name"`
	Default              bool              `yaml:"default"`
	ProjectsForDirectory map[string]string `yaml:"projects_for_directory"`
}

func (c EvergreenYamlConfiguration) getProjectForDirectory(directory string) string {
	for k, v := range c.ProjectsForDirectory {
		if k == directory {
			fmt.Printf("Current directory is associated with project [%s], updating active project.\n", v)
			return v
		}
	}

	for _, p := range c.Projects {
		if p.Default {
			fmt.Printf("Default project found [%s].\n", p.Name)
			return p.Name
		}
	}
	return ""
}

func loadConfigFromEvergreenYaml(fileBytes []byte) (Config, error) {
	evgConfig := EvergreenYamlConfiguration{}
	if err := yaml.Unmarshal(fileBytes, &evgConfig); err != nil {
		return Config{}, err
	}

	currentDirectory, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}

	return Config{
		User:           evgConfig.User,
		BaseUrl:        evgConfig.BaseUrl,
		ApiKey:         evgConfig.ApiKey,
		DefaultProject: evgConfig.getProjectForDirectory(currentDirectory),
	}, nil
}
