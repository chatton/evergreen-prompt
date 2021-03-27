package client

import (
	"chatton.com/evergreen-prompt/pkg/evergreen/project"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

type EvergreenClient struct {
	apiKey   string
	username string
	baseUrl  string
	client   *http.Client

	projects       []string
	currentProject string
	ActiveProject  string
}

type Config struct {
	User    string `json:"user"`
	BaseUrl string `json:"baseUrl"`
	ApiKey  string `json:"apiKey"`
}

// LoadConfig loads the config file that contains the user, baseurl and apikey
// required to interact with the evergreen api.
func LoadConfig() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	configFilePath := path.Join(home, ".evergreen-prompt.json")
	bytes, err := ioutil.ReadFile(configFilePath)
	config := Config{}
	if err := json.Unmarshal(bytes, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func NewEvergreenClient() (*EvergreenClient, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	return newEvergreenClientFromConfig(config), nil
}

func newEvergreenClientFromConfig(config Config) *EvergreenClient {
	return &EvergreenClient{
		apiKey:   config.ApiKey,
		username: config.User,
		baseUrl:  config.BaseUrl,
		client:   &http.Client{},
		projects: nil,
	}
}

func (c *EvergreenClient) get(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.baseUrl, endpoint), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Api-User", c.username)
	req.Header.Set("Api-Key", c.apiKey)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(res.Body)
}

func (c *EvergreenClient) GetProjects() ([]string, error) {
	if c.projects != nil {
		return c.projects, nil
	}
	b, err := c.get("rest/v1/projects")
	if err != nil {
		return nil, err
	}
	allProjects := &project.Response{}
	if err := json.Unmarshal(b, allProjects); err != nil {
		return nil, err
	}

	c.projects = allProjects.Projects
	return allProjects.Projects, nil

}
