package tohyo_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shota-imoto/tohyo"
)

func TestInteractiveVote(t *testing.T) {
	t.Setenv("CANDIDATE_PATH", "testdata/candidate.json")
	tests := []struct {
		name     string
		testPath string
		wants    tohyo.ResultsMap
	}{
		{
			name:     "two_voter",
			testPath: "./testdata/wants/two_voter.txt",
			wants:    map[string]int{"gimuki": 4, "kamuki": 2, "chosonmuki": 0},
		},
		{
			name:     "one_voter",
			testPath: "./testdata/wants/one_voter.txt",
			wants:    map[string]int{"gimuki": 2, "kamuki": 1, "chosonmuki": 0},
		},
		{
			name:     "same_candidate",
			testPath: "./testdata/wants/same_candidate.txt",
			wants:    map[string]int{"gimuki": 0, "kamuki": 2, "chosonmuki": 0},
		},
		{
			name:     "not_exist_candidate",
			testPath: "./testdata/wants/not_exist_candidate.txt",
			wants:    map[string]int{"gimuki": 0, "kamuki": 0, "chosonmuki": 0},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r, err := os.Open(tt.testPath)

			if err != nil {
				t.Fatalf("TestInteractiveVote: %v", err)
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
