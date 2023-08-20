package config

import "github.com/caarlos0/env/v6"

type Config struct {
	CandidatePath string `env:"CANDIDATE_PATH" envDefault:"candidate.json"`
	RightPath     string `env:"RIGHT_PATH" envDefault:"right.json"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
