package countvotes

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/polls/repository"
)

type rankedVoteNode struct {
	NextPreferences map[int]*rankedVoteNode
	Votes           int
}

func newRankedVoteNode() *rankedVoteNode {
	return &rankedVoteNode{
		NextPreferences: make(map[int]*rankedVoteNode, 0),
	}
}

func countRankedChoice(poll model.Poll) {
	t, err := loadInitialVoteTree(poll.PollId)
	if err != nil {
		log.Printf("failed fetching vote items: %s\n", err.Error())
		return
	}

	target := t.Votes/2 + 1
	result := model.RankedChoicePollResult{
		PollId:        poll.PollId,
		Discriminator: model.DiscriminatorResult,
		Votes:         []model.RankedChoiceOptionVoteCount{},
	}
	for _, o := range poll.Options {
		result.Votes = append(result.Votes, model.RankedChoiceOptionVoteCount{
			Option:     o,
			RoundVotes: []int{},
		})
	}

	eliminated := make([]int, 0)
	for {
		mostVotes := -1
		leastVotes := math.MaxInt32
		leastPopular := []int{}
		for i := range result.Votes {
			if isEliminated(i, eliminated) {
				continue
			}

			if t.NextPreferences[i] == nil {
				t.NextPreferences[i] = newRankedVoteNode()
			}

			result.Votes[i].RoundVotes = append(result.Votes[i].RoundVotes, t.NextPreferences[i].Votes)

			if t.NextPreferences[i].Votes > mostVotes {
				mostVotes = t.NextPreferences[i].Votes
			}
			if t.NextPreferences[i].Votes < leastVotes {
				leastVotes = t.NextPreferences[i].Votes
				leastPopular = []int{i}
			} else if t.NextPreferences[i].Votes == leastVotes {
				leastPopular = append(leastPopular, i)
			}
		}

		if mostVotes >= target || len(t.NextPreferences) <= 1 {
			break
		}

		for _, i := range leastPopular {
			t.remove(i)
		}

		eliminated = append(eliminated, leastPopular...)
	}

	err = repository.InsertNewPollItem(result)
	if err != nil {
		log.Printf("failed to write result to DB %s\n", err.Error())
	}
}

func loadInitialVoteTree(pollId string) (*rankedVoteNode, error) {
	t := newRankedVoteNode()
	var continuation *string = nil
	for {
		page := make([]model.RankedChoiceVote, 0)
		continuation, err := repository.GetPollVoteItems(pollId, continuation, &page)
		if err != nil {
			return nil, err
		}

		for _, v := range page {
			t.addVote(v.Ranked)
		}

		if continuation == nil {
			return t, nil
		}
	}
}

func isEliminated(o int, eliminated []int) bool {
	for _, e := range eliminated {
		if e == o {
			return true
		}
	}

	return false
}

func (t *rankedVoteNode) addVote(ranked []int) {
	t.Votes++
	if len(ranked) == 0 {
		return
	}

	n := t.NextPreferences[ranked[0]]
	if n == nil {
		n = newRankedVoteNode()
		t.NextPreferences[ranked[0]] = n
	}

	n.addVote(ranked[1:])
}

func (t *rankedVoteNode) remove(o int) {
	for co, c := range t.NextPreferences {
		if co != o {
			c.remove(o)
		}
	}

	if r := t.NextPreferences[o]; r != nil {
		for co, c := range r.NextPreferences {
			n := t.NextPreferences[co]
			if n == nil {
				n = newRankedVoteNode()
				t.NextPreferences[co] = n
			}
			n.merge(c)
		}
	}

	delete(t.NextPreferences, o)
}

func (t *rankedVoteNode) merge(other *rankedVoteNode) {
	t.Votes += other.Votes
	for o, c := range other.NextPreferences {
		n := t.NextPreferences[o]
		if n == nil {
			n = newRankedVoteNode()
			t.NextPreferences[o] = n
		}
		n.merge(c)
	}
}

func (t *rankedVoteNode) writeParts() []string {
	parts := make([]string, 0)
	parts = append(parts, fmt.Sprintf("%d;%d;", t.Votes, len(t.NextPreferences)))
	for k, t := range t.NextPreferences {
		parts = append(parts, fmt.Sprintf("%d;", k))
		parts = append(parts, t.writeParts()...)
	}

	return parts
}

func (t *rankedVoteNode) write() string {
	return strings.Join(t.writeParts(), "")
}

type RankedNodeTreeParser struct {
	ser string
}

func (p *RankedNodeTreeParser) nextToken() (int, error) {
	c := 0
	for c < len(p.ser) && p.ser[c] != ';' {
		c++
	}
	if c >= len(p.ser) {
		return 0, errors.New("unexpected end of stream")
	}

	val, err := strconv.Atoi(p.ser[:c])
	if err != nil {
		return 0, err
	}

	p.ser = p.ser[c+1:]
	return val, nil
}

func (p *RankedNodeTreeParser) parse() (*rankedVoteNode, error) {
	votes, err := p.nextToken()
	if err != nil {
		return nil, err
	}

	npc, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	nextPreferences := make(map[int]*rankedVoteNode, 0)
	for i := 0; i < npc; i++ {
		o, err := p.nextToken()
		if err != nil {
			return nil, err
		}
		n, err := p.parse()
		if err != nil {
			return nil, err
		}
		nextPreferences[o] = n
	}

	return &rankedVoteNode{
		Votes:           votes,
		NextPreferences: nextPreferences,
	}, nil
}
