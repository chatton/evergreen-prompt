package client

import (
	"chatton.com/evergreen-prompt/pkg/evergreen/project"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type EvergreenClient struct {
	apiKey   string
	username string
	baseUrl  string
	client   *http.Client

	projects       []string
	currentProject string
}

func NewEvergreenClient(username, apiKey, baseUrl string) *EvergreenClient {
	return &EvergreenClient{
		apiKey:   apiKey,
		username: username,
		baseUrl:  baseUrl,
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
