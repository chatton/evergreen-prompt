package evgprompt

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

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
