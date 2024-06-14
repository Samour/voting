package castvote

import (
	"errors"
	"fmt"
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

	return &castVoteModel{
		Poll:    poll,
		MayVote: poll.Status == "voting",
		Voted:   -1,
	}, nil
}

func castVote(pollId string, option int) (*castVoteModel, error) {
	poll, err := repository.GetPollItem(pollId)
	if err != nil {
		return nil, err
	}

	if option < 0 || option >= len(poll.Options) {
		return nil, errors.New("invalid option provided")
	}

	voteId := utils.IdGen()
	discriminator := fmt.Sprintf("vote:%s", voteId)
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
