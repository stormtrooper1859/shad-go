package firewall

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type RulesConfig struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Endpoint               string   `yaml:"endpoint"`
	ForbiddenUserAgents    []string `yaml:"forbidden_user_agents"`
	ForbiddenHeaders       []string `yaml:"forbidden_headers"`
	RequiredHeaders        []string `yaml:"required_headers"`
	MaxRequestLengthBytes  int      `yaml:"max_request_length_bytes"`
	MaxResponseLengthBytes int      `yaml:"max_response_length_bytes"`
	ForbiddenResponseCodes []int    `yaml:"forbidden_response_codes"`
	ForbiddenRequestRe     []string `yaml:"forbidden_request_re"`
	ForbiddenResponseRe    []string `yaml:"forbidden_response_re"`
}

func GetConfig(path string) (*RulesConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config RulesConfig

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
