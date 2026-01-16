package parser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Project Project `yaml:"project"`
	// Feature     Feature  `yaml:"feature"`
	Entities    []string `yaml:"entities"`
	CustomLogic []string `yaml:"custom_logic"`
}

type Project struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Port     int    `yaml:"port"`
	Location string `yaml:"location"`
	Database string `yaml:"db"`
	Router   string `yaml:"router"`
}

// type Feature struct {
// 	Database FeatureItem `yaml:"database"`
// 	Cache    FeatureItem `yaml:"cache"`
// 	Queue    FeatureItem `yaml:"queue"`
// 	Auth     FeatureItem `yaml:"auth"`
// }

type FeatureItem struct {
	Name string `yaml:"name"`
}

type ProjectYAML struct {
	Name     string   `json:"name" yaml:"name"`
	Layers   []string `json:"layers" yaml:"layers"`
	Features []string `json:"features" yaml:"features"`
}

func ReadYAML(yamlPath string) (*Config, error) {
	yamlByte, err := os.ReadFile(yamlPath)
	if err != nil {
		return &Config{}, err
	}

	var yamlProject Config

	err = yaml.Unmarshal(yamlByte, &yamlProject)
	if err != nil {
		fmt.Println("error: ", err)
		return &Config{}, err
	}

	return &yamlProject, nil
}
