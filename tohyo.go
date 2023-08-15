package tohyo

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type ResultsMap map[string]int

type Vote struct {
	Candidates Candidates
	Voters     int
	Right      int
	ResultsMap ResultsMap
	r          io.Reader
}

var right int = 2

func NewVote(r io.Reader) (Vote, error) {
	candidates, err := LoadCandidates()
	if err != nil {
		return Vote{}, fmt.Errorf("NewVote: %w", err)
	}
	rm := candidates.NewResultsMap()

	return Vote{Candidates: candidates, Voters: 0, Right: right, ResultsMap: rm, r: r}, nil
}

func InteractiveVote(r io.Reader) error {
	fmt.Println("投票人数を入力してください")
	v, err := NewVote(r)
	if err != nil {
		return fmt.Errorf("InteractiveVote: %w", err)
	}

	s := bufio.NewScanner(v.r)
	s.Scan()
	input := s.Text()
	i, ok := strconv.Atoi(input)
	if ok != nil {
		fmt.Println("数字以外が入力されたので処理を終了します")
		return nil
	}
	v.Voters = i
	fmt.Printf("人数テスト: %v", v)

	fmt.Println("投票を開始します")
	v.Start()
	fmt.Println(v.ResultsMap)

	return nil
}

func (v Vote) Start() error {
	for i := 0; i < v.Voters; i++ {
		fmt.Printf(
			"投票権は「%d」与えられます。いいと思った順に候補を入力してEnterを押してください。\n",
			v.Right,
		)
		fmt.Println("※投票順に重み付けされます")
		s := bufio.NewScanner(v.r)

		// 同じ候補に2回投票しないために、投票済み候補を保持するslice
		voted := make([]string, v.Right)

		for j := 0; j < v.Right; j++ {
			fmt.Printf("%d回目の投票権を行使します。\n", j+1)
			fmt.Printf("候補を入力してEnterを押してください: %v\n", v.CandidatesString())
			for s.Scan() {
				t := s.Text()

				if slices.Contains(voted, t) {
					fmt.Println("同じ候補に2回以上投票できません。未投票の候補を入力してください")
					fmt.Printf("入力済み: %s\n", strings.Join(voted, ","))
					continue
				}

				err := v.Count(t, j)

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

var ErrNotCandidate = errors.New("not candidate")

func (v Vote) Count(str string, priority int) error {
	if _, ok := v.ResultsMap[str]; !ok {
		return ErrNotCandidate
	}
	v.ResultsMap[str] += v.Right - priority
	return nil
}
