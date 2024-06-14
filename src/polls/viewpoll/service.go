package viewpoll

import (
	"errors"
	"time"

	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/polls/repository"
)

func getPoll(id string) (*model.ViewPollModel, error) {
	poll, err := repository.GetPollItem(id)
	if err != nil {
		return nil, err
	}

	return ToViewPollModel(poll), nil
}

func ToViewPollModel(p *model.Poll) *model.ViewPollModel {
	return &model.ViewPollModel{
		Poll:            p,
		OobStatusUpdate: false,
	}
}

func updateStatus(id string, status string) (*model.ViewPollModel, error) {
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
	} else if status == "closed" {
		if poll.Status != "voting" {
			return nil, errors.New("voting is not currently open on poll")
		}
		poll.Statistics.ClosedAt = time.Now().In(time.UTC).Format(time.RFC3339)
	} else {
		return nil, errors.New("unknown status")
	}

	poll.Status = status
	err = repository.UpdatePollItem(poll)
	if err != nil {
		return nil, err
	}

	return &model.ViewPollModel{
		Poll:            poll,
		OobStatusUpdate: true,
	}, nil
}
