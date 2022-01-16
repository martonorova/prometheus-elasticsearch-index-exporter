package config

import (
	"fmt"
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Global        GlobalConfig        `yaml:",omitempty"`
	Elasticsearch ElasticsearchConfig `yaml:",omitempty"`
	AllMetrics    MetricsConfig       `yaml:"-"`
	Server        ServerConfig        `yaml:",omitempty"`
}

type GlobalConfig struct {
	ConfigVersion   int           `yaml:"config_version,omitempty"`
	RefreshInterval time.Duration `yaml:"refresh_interval,omitempty"` // implicitly parsed with time.ParseDuration()
}

type ElasticsearchConfig struct {
	Hosts []string `yaml:",omitempty"`
}

type MetricConfig struct {
	Type       string              `yaml:",omitempty"`
	Name       string              `yaml:",omitempty"`
	Help       string              `yaml:",omitempty"`
	Match      string              `yaml:",omitempty"`
	Retention  time.Duration       `yaml:",omitempty"` // implicitly parsed with time.ParseDuration()
	Value      string              `yaml:",omitempty"`
	Cumulative bool                `yaml:",omitempty"`
	Buckets    []float64           `yaml:",flow,omitempty"`
	Quantiles  map[float64]float64 `yaml:",flow,omitempty"`
	MaxAge     time.Duration       `yaml:"max_age,omitempty"`
	Labels     map[string]string   `yaml:",omitempty"`
	// LabelTemplates       []template.Template `yaml:"-"` // parsed version of Labels, will not be serialized to yaml.
	// ValueTemplate        template.Template   `yaml:"-"` // parsed version of Value, will not be serialized to yaml.
	// DeleteMatch          string              `yaml:"delete_match,omitempty"`
	// DeleteLabels         map[string]string   `yaml:"delete_labels,omitempty"` // TODO: Make sure that DeleteMatch is not nil if DeleteLabels are used.
	// DeleteLabelTemplates []template.Template `yaml:"-"`                       // parsed version of DeleteLabels, will not be serialized to yaml.
}

type MetricsConfig []MetricConfig

type ServerConfig struct {
	Protocol   string `yaml:",omitempty"`
	Host       string `yaml:",omitempty"`
	Port       int    `yaml:",omitempty"`
	Path       string `yaml:",omitempty"`
	Cert       string `yaml:",omitempty"`
	Key        string `yaml:",omitempty"`
	ClientCA   string `yaml:"client_ca,omitempty"`
	ClientAuth string `yaml:"client_auth,omitempty"`
}

func LoadConfigFile(filename string) (*Config, string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to load %v: %v", filename, err.Error())
	}
	cfg, warn, err := LoadConfigString(content)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to load %v: %v", filename, err.Error())
	}
	return cfg, warn, err
}

func LoadConfigString(content []byte) (*Config, string, error) {
	// TODO, handle config version here
	cfg, err := unmarshal(content)
	return cfg, "", err
}

func unmarshal(content []byte) (*Config, error) {
	cfg := &Config{}
	err := yaml.Unmarshal(content, cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %v. make sure to use 'single quotes' around strings with special characters (like match patterns or label templates), and make sure to use '-' only for lists (metrics) but not for maps (labels)", err.Error())
	}
	// TODO
	// err = AddDefaultsAndValidate(cfg)
	// if err != nil {
	// 	return nil, err
	// }
	return cfg, nil
}
