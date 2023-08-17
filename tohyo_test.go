package tohyo_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shota-imoto/tohyo"
)

func TestInteractiveVote(t *testing.T) {
	tests := []struct {
		name     string
		testPath string
		wants    tohyo.ResultsMap
	}{
		{
			name:     "two_voter",
			testPath: "./testdata/two_voter.txt",
			wants:    map[string]int{"buggy": 4, "jungle": 2, "cave": 0},
		},
		{
			name:     "one_voter",
			testPath: "./testdata/one_voter.txt",
			wants:    map[string]int{"buggy": 2, "jungle": 1, "cave": 0},
		},
		{
			name:     "same_candidate",
			testPath: "./testdata/same_candidate.txt",
			wants:    map[string]int{"buggy": 0, "jungle": 2, "cave": 0},
		},
		{
			name:     "not_exist_candidate",
			testPath: "./testdata/not_exist_candidate.txt",
			wants:    map[string]int{"buggy": 0, "jungle": 0, "cave": 0},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r, err := os.Open(tt.testPath)

			if err != nil {
				t.Errorf("TestInteractiveVote: %v", err)
			}
			defer r.Close()

			v, err := tohyo.InteractiveVote(r)
			if err != nil {
				t.Errorf("TestInteractiveVote: %v", err)
			}
			if diff := cmp.Diff(v.ResultsMap, tt.wants); diff != "" {
				t.Errorf("Result is mismatch (-ResultMap +ResultMap): %s\n", diff)
			}
		})
	}
}
