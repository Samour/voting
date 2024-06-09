package polls

import (
	"fmt"
	"time"

	"github.com/Samour/voting/utils"
)

func CastVote(pollId string, option int) (*Poll, error) {
	voteId := utils.IdGen()
	discriminator := fmt.Sprintf("vote:%s", voteId)
	vote := Vote{
		PollId:        pollId,
		Discriminator: discriminator,
		Option:        option,
		CastAt:        time.Now().In(time.UTC).Format(time.RFC3339),
	}

	err := recordVote(&vote)
	if err != nil {
		return nil, err
	}

	return getPollItem(pollId)
}
