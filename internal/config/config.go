package internal

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

// Read read the config file from the configFileName file at home dir
func Read() Config {
	cfgFilePath, err := getConfigFilePath()
	if err != nil {
		fmt.Println(err)
		return Config{}

	}
	data, err := os.ReadFile(cfgFilePath)
	if err != nil {
		fmt.Println(err)
		return Config{}
	}
	cfg := Config{}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		fmt.Println(err)
		return Config{}
	}

	return cfg
}

// SetUser set the current username of the config and also write/update  the configFileName file at home dir
func (c Config) SetUser(userName string) error {
	c.CurrentUserName = userName
	return write(c)

}

func getConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", nil
	}

	return homedir + "/" + configFileName, nil
}
func write(cfg Config) error {
	cfgFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(cfgFilePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	// data, err := json.Marshal(cfg)
	// lets use mashallindent to indent the json
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	// os.WriteFile(cfgFilePath, , fs.FileMode(os.O_APPEND))
	_, err = file.Write(data)
	return err
}
