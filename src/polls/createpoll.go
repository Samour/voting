package polls

import (
	"github.com/Samour/voting/utils"
)

func CreatePoll() (*string, error) {
	id := utils.IdGen()
	poll := Poll{
		PollId:        id,
		Discriminator: "poll",
		Status:        "draft",
		Name:          "",
		Options:       []string{""},
	}

	err := insertNewPollItem(&poll)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
