package tohyo

import (
	"encoding/json"
	"os"

	"github.com/shota-imoto/tohyo/config"
)

type Rights []int

func LoadRights() (Rights, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(cfg.RightPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var loaded struct {
		Rights Rights `json:"rights"`
	}

	if err := json.NewDecoder(f).Decode(&loaded); err != nil {
		return nil, err
	}
	return loaded.Rights, nil
}
