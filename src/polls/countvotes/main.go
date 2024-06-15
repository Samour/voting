package countvotes

import (
	"log"

	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/polls/repository"
)

func CountVotes(pollId string) {
	poll := &model.Poll{}
	err := repository.GetPollItem(pollId, model.DiscriminatorPoll, poll)
	if err != nil {
		log.Printf("failed fetching poll: %s\n", err.Error())
		return
	}

	if poll.AggregationType == model.PollAggregationTypeFirstPastThePost {
		countFptp(poll)
	} else if poll.AggregationType == model.PollAggregationTypeRankedChoice {
		countRankedChoice(poll)
	}
}
