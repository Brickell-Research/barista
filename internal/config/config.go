package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Service struct {
	Name         string `yaml:"name"`
	URL          string `yaml:"url"`
	ProviderName string
}

func (s Service) Key() string {
	return s.ProviderName + "/" + s.Name
}

type Provider struct {
	Name     string    `yaml:"name"`
	Services []Service `yaml:"services"`
}

type Config struct {
	Providers  []Provider `yaml:"providers"`
	OutputDir  string     `yaml:"output_dir"`
	OutputRepo string     `yaml:"output_repo"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if cfg.OutputDir == "" {
		cfg.OutputDir = "output"
	}

	for i := range cfg.Providers {
		for j := range cfg.Providers[i].Services {
			cfg.Providers[i].Services[j].ProviderName = cfg.Providers[i].Name
		}
	}

	return &cfg, nil
}

func (c *Config) AllServices() []Service {
	var services []Service
	for _, p := range c.Providers {
		services = append(services, p.Services...)
	}
	return services
}

func (c *Config) FindService(key string) (*Service, bool) {
	for _, p := range c.Providers {
		for _, s := range p.Services {
			if s.Key() == key {
				return &s, true
			}
		}
	}
	return nil, false
}
