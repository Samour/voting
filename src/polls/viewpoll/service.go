package viewpoll

import (
	"errors"
	"time"

	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/polls/repository"
)

func updateStatus(id string, status string) (*model.Poll, error) {
	poll, err := repository.GetPollItem(id)
	if err != nil {
		return nil, err
	}
	if poll == nil {
		return nil, nil
	}

	if status == "voting" {
		if poll.Status != "draft" {
			return nil, errors.New("cannot open voting on poll")
		}
		poll.Statistics.OpenedAt = time.Now().In(time.UTC).Format(time.RFC3339)
	} else {
		return nil, errors.New("unknown status")
	}

	poll.Status = status
	err = repository.UpdatePollItem(poll)
	if err != nil {
		return nil, err
	}

	return poll, nil
}