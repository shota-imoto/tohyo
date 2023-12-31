package tohyo

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/shota-imoto/tohyo/config"
)

type Candidate struct {
	Name string `json:"name"`
}

type Candidates []Candidate

func LoadCandidates() (Candidates, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(cfg.CandidatePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var loaded struct {
		Candidates Candidates `json:"candidates"`
	}

	if err := json.NewDecoder(f).Decode(&loaded); err != nil {
		return nil, err
	}
	return loaded.Candidates, nil
}

func (cs Candidates) String() string {
	list := make([]string, len(cs))

	for i, c := range cs {
		list[i] = c.Name
	}
	return strings.Join(list, ",")
}

func (cs Candidates) NewResultsMap() ResultsMap {
	rm := make(ResultsMap)

	for _, c := range cs {
		rm[c.Name] = 0
	}
	return rm
}
