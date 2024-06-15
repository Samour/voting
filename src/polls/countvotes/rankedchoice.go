package countvotes

import (
	"log"
	"math"

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

func countRankedChoice(poll *model.Poll) {
	t := make(map[int]*rankedVoteNode, 0)
	totalVotes, err := loadInitialVoteTree(poll.PollId, t)
	if err != nil {
		log.Printf("failed fetching vote items: %s\n", err.Error())
		return
	}

	target := totalVotes/2 + 1
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

			if t[i] == nil {
				t[i] = newRankedVoteNode()
			}

			result.Votes[i].RoundVotes = append(result.Votes[i].RoundVotes, t[i].Votes)

			if t[i].Votes > mostVotes {
				mostVotes = t[i].Votes
			}
			if t[i].Votes < leastVotes {
				leastVotes = t[i].Votes
				leastPopular = []int{i}
			} else if t[i].Votes == leastVotes {
				leastPopular = append(leastPopular, i)
			}
		}

		if mostVotes >= target || len(t) <= 1 {
			break
		}

		for _, i := range leastPopular {
			for o, f := range t {
				if o != i {
					f.remove(i)
				}
			}

			for o, n := range t[i].NextPreferences {
				t[o].merge(n)
			}
			delete(t, i)
		}

		eliminated = append(eliminated, leastPopular...)
	}

	err = repository.InsertNewPollItem(result)
	if err != nil {
		log.Printf("failed to write result to DB %s\n", err.Error())
	}
}

func loadInitialVoteTree(pollId string, t map[int]*rankedVoteNode) (int, error) {
	totalVotes := 0
	var continuation *string = nil
	for {
		page := make([]model.RankedChoiceVote, 0)
		continuation, err := repository.GetPollVoteItems(pollId, continuation, &page)
		if err != nil {
			return -1, err
		}

		for _, v := range page {
			totalVotes++
			n := t[v.Ranked[0]]
			if n == nil {
				n = newRankedVoteNode()
				t[v.Ranked[0]] = n
			}
			n.addVote(v.Ranked[1:])
		}

		if continuation == nil {
			return totalVotes, nil
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
