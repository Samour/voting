package editpoll

import (
	"errors"

	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/polls/repository"
)

func getPoll(id string) (*editPollModel, error) {
	poll, err := repository.GetPollItem(id)
	if err != nil {
		return nil, err
	}
	if poll == nil {
		return nil, nil
	}

	return &editPollModel{
		Poll:    poll,
		MayEdit: poll.Status == "draft",
	}, nil
}

func updatePollDetails(id string, d pollDetails) (*model.Poll, error) {
	poll, err := repository.GetPollItem(id)
	if err != nil {
		return nil, err
	}
	if poll == nil {
		return nil, nil
	}

	if poll.Status != "draft" {
		return nil, errors.New("cannot edit poll that is not in draft status")
	}

	poll.Name = d.Name
	poll.Options = d.Options
	err = repository.UpdatePollItem(poll)
	if err != nil {
		return nil, err
	}

	return poll, nil
}

func patchPollOptions(options []string, u pollOptionsUpdate) []string {
	if u.Remove >= 0 && u.Remove < len(options) {
		options = append(options[:u.Remove], options[u.Remove+1:]...)
	}
	if u.Add || len(options) == 0 {
		options = append(options, "")
	}

	return options
}
