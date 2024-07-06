package createpoll

import (
	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/polls/repository"
	"github.com/Samour/voting/utils"
)

func createPoll() (string, error) {
	id := utils.IdGen()
	poll := model.Poll{
		PollId:          id,
		Discriminator:   model.DiscriminatorPoll,
		Status:          model.PollStatusDraft,
		AggregationType: model.PollAggregationTypeFirstPastThePost,
		Name:            "",
		Options:         []string{""},
	}

	err := repository.InsertNewPollItem(&poll)
	if err != nil {
		return "", err
	}

	return id, nil
}
