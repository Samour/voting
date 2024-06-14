package castvote

import (
	"fmt"
	"time"

	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/polls/repository"
	"github.com/Samour/voting/utils"
)

func castVote(pollId string, option int) (*model.Poll, error) {
	voteId := utils.IdGen()
	discriminator := fmt.Sprintf("vote:%s", voteId)
	vote := model.Vote{
		PollId:        pollId,
		Discriminator: discriminator,
		Option:        option,
		CastAt:        time.Now().In(time.UTC).Format(time.RFC3339),
	}

	err := repository.RecordVote(&vote)
	if err != nil {
		return nil, err
	}

	return repository.GetPollItem(pollId)
}
