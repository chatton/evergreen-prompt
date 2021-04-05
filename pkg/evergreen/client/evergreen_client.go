package client

import (
	"bytes"
	"chatton.com/evergreen-prompt/pkg/evergreen/patch"
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

	currentPatches []patch.Patch
	DefaultProject string
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
		apiKey:         config.ApiKey,
		username:       config.User,
		baseUrl:        config.BaseUrl,
		client:         &http.Client{},
		projects:       nil,
		DefaultProject: config.DefaultProject,
	}
}

func (c *EvergreenClient) get(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.baseUrl, endpoint), nil)
	if err != nil {
		return nil, err
	}
	if req.Body != nil {
		defer req.Body.Close()
	}

	c.setHeaders(req)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(res.Body)
}

func (c *EvergreenClient) patch(endpoint string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%s", c.baseUrl, endpoint), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	if req.Body != nil {
		defer req.Body.Close()
	}

	c.setHeaders(req)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(res.Body)
}

func (c *EvergreenClient) setHeaders(req *http.Request) {
	req.Header.Set("Api-User", c.username)
	req.Header.Set("Api-Key", c.apiKey)
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

func (c *EvergreenClient) PatchPatch(patchId string, body patch.Body) ([]string, error) {
	bytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	_, err = c.patch(fmt.Sprintf("rest/v2/patches/%s", patchId), bytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *EvergreenClient) GetPatches() ([]patch.Patch, error) {
	if c.currentPatches != nil {
		return c.currentPatches, nil
	}

	b, err := c.get(fmt.Sprintf("rest/v2/users/%s/patches", c.username))
	if err != nil {
		return []patch.Patch{}, err
	}

	var result []patch.Patch

	var patches []patch.Patch
	if err := json.Unmarshal(b, &patches); err != nil {
		return []patch.Patch{}, err
	}

	// only return the patches that are for the active project.
	for _, p := range patches {
		if p.ProjectId == c.DefaultProject {
			result = append(result, p)
		}
	}

	c.currentPatches = result
	return result, err
}
