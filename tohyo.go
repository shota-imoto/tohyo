package tohyo

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type ResultsMap map[string]int

type Vote struct {
	Candidates Candidates
	Voters     int
	Rights     []int
	ResultsMap ResultsMap
	s          *bufio.Scanner
}

func NewVote(r io.Reader) (Vote, error) {
	candidates, err := LoadCandidates()
	if err != nil {
		return Vote{}, fmt.Errorf("NewVote: %w", err)
	}

	rm := candidates.NewResultsMap()
	s := bufio.NewScanner(r)

	rights, err := LoadRights()
	if err != nil {
		return Vote{}, fmt.Errorf("NewVote: %w", err)
	}

	return Vote{Candidates: candidates, Voters: 0, Rights: rights, ResultsMap: rm, s: s}, nil
}

func (v Vote) Scan() bool {
	return v.s.Scan()
}

func (v Vote) Text() string {
	return v.s.Text()
}

func InteractiveVote(r io.Reader) (Vote, error) {
	fmt.Println("投票人数を入力してください")
	v, err := NewVote(r)
	if err != nil {
		return Vote{}, fmt.Errorf("InteractiveVote: %w", err)
	}

	v.Scan()
	input := v.Text()

	i, ok := strconv.Atoi(input)
	if ok != nil {
		fmt.Println("数字以外が入力されたので処理を終了します")
		return v, nil
	}
	v.Voters = i
	fmt.Printf("人数テスト: %v", v)

	fmt.Println("投票を開始します")
	v.Start()

	fmt.Printf("結果:\n%s", v.ResultString())

	return v, nil
}

func (v Vote) Start() error {
	for i := 0; i < v.Voters; i++ {
		fmt.Printf(
			"投票権は「%d」与えられます。いいと思った順に候補を入力してEnterを押してください。\n",
			len(v.Rights),
		)
		fmt.Println("※投票順に重み付けされます")

		// 同じ候補に2回投票しないために、投票済み候補を保持するslice
		voted := make([]string, len(v.Rights))

		for j, r := range v.Rights {
			fmt.Printf("%d回目の投票権を行使します。\n", j+1)
			fmt.Printf("候補を入力してEnterを押してください: %v\n", v.CandidatesString())

			for v.Scan() {
				t := v.Text()

				if slices.Contains(voted, t) {
					fmt.Println("同じ候補に2回以上投票できません。未投票の候補を入力してください")
					fmt.Printf("入力済み: %s\n", strings.Join(voted, ","))
					continue
				}

				err := v.Count(t, r)

				if err != nil {
					if errors.Is(ErrNotCandidate, err) {
						fmt.Println("候補にない名前が入力されました。候補の中から入力してください")
						continue
					}
					return err
				}
				voted[j] = t
				break
			}
		}
	}
	return nil
}

func (v Vote) CandidatesString() string {
	return v.Candidates.String()
}

func (v Vote) ResultString() string {
	type SortedResult struct {
		Key   string
		Value int
	}

	sortedResults := make([]SortedResult, len(v.ResultsMap))

	i := 0
	for k, v := range v.ResultsMap {
		sortedResults[i] = SortedResult{Key: k, Value: v}
		i++
	}
	sort.Slice(sortedResults, func(i, j int) bool { return sortedResults[i].Value > sortedResults[j].Value })

	str := ""

	for i, r := range sortedResults {
		str += fmt.Sprintf("%d位: %s -> %d\n", i+1, r.Key, r.Value)
	}

	return str
}

var ErrNotCandidate = errors.New("not candidate")

func (v Vote) Count(str string, weight int) error {
	if _, ok := v.ResultsMap[str]; !ok {
		return ErrNotCandidate
	}
	v.ResultsMap[str] += weight
	return nil
}
