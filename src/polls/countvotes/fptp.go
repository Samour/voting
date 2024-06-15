package countvotes

import (
	"log"

	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/polls/repository"
)

func countFptp(poll *model.Poll) {
	voteCounts := make(map[int]int, 0)
	var continuation *string = nil

	for {
		page, err := repository.GetPollVoteItems(poll.PollId, continuation)
		if err != nil {
			log.Printf("failed fetching vote items: %s\n", err.Error())
			return
		}

		for _, vote := range page.Items {
			voteCounts[vote.Option] = voteCounts[vote.Option] + 1
		}

		continuation = page.LastEvaluatedKey
		if continuation == nil {
			break
		}
	}

	counts := make([]model.OptionVoteCount, len(poll.Options))
	for i, o := range poll.Options {
		counts[i] = model.OptionVoteCount{
			Option:    o,
			VoteCount: voteCounts[i],
		}
	}

	result := &model.PollResult{
		PollId:        poll.PollId,
		Discriminator: model.DiscriminatorResult,
		Votes:         counts,
	}
	err := repository.InsertNewPollItem(result)
	if err != nil {
		log.Printf("failed to insert poll result: %s\n", err.Error())
	}
}
