package castvote

import (
	"errors"
	"time"

	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/polls/repository"
	"github.com/Samour/voting/utils"
)

func getPoll(id string) (*castVoteModel, error) {
	poll, err := repository.GetPollItem(id)
	if err != nil {
		return nil, err
	}
	if poll == nil {
		return nil, nil
	}

	var rco *rankedChoiceOptions = nil
	if poll.AggregationType == model.PollAggregationTypeRankedChoice {
		uo := make([]voteOption, len(poll.Options))
		for i, o := range poll.Options {
			uo[i] = voteOption{
				Index:  i,
				Option: o,
			}
		}
		rco = &rankedChoiceOptions{
			Unselected: uo,
		}
	}

	return &castVoteModel{
		Poll:    poll,
		Rco:     rco,
		MayVote: poll.Status == model.PollStatusVoting,
		Voted:   -1,
	}, nil
}

func castVote(pollId string, option int) (*castVoteModel, error) {
	poll, err := repository.GetPollItem(pollId)
	if err != nil {
		return nil, err
	}

	if poll.Status != model.PollStatusVoting {
		return nil, errors.New("poll is not open for voting")
	}

	if option < 0 || option >= len(poll.Options) {
		return nil, errors.New("invalid option provided")
	}

	voteId := utils.IdGen()
	discriminator := model.DiscriminatorVote + voteId
	vote := model.Vote{
		PollId:        pollId,
		Discriminator: discriminator,
		Option:        option,
		CastAt:        time.Now().In(time.UTC).Format(time.RFC3339),
	}

	err = repository.RecordVote(&vote)
	if err != nil {
		return nil, err
	}

	return &castVoteModel{
		Poll:    poll,
		MayVote: true,
		Voted:   option,
	}, nil
}
