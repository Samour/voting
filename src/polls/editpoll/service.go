package editpoll

import (
	"errors"

	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/polls/repository"
	"github.com/Samour/voting/polls/viewpoll"
)

func getPoll(id string) (*editPollModel, error) {
	poll := &model.Poll{}
	err := repository.GetPollItem(id, model.DiscriminatorPoll, poll)
	if err != nil {
		return nil, err
	}
	if len(poll.PollId) == 0 {
		return nil, nil
	}

	return &editPollModel{
		Poll:    poll,
		MayEdit: poll.Status == model.PollStatusDraft,
	}, nil
}

func updatePollDetails(id string, d pollDetails) (*model.ViewPollModel, error) {
	poll := &model.Poll{}
	err := repository.GetPollItem(id, model.DiscriminatorPoll, poll)
	if err != nil {
		return nil, err
	}
	if len(poll.PollId) == 0 {
		return nil, nil
	}

	if poll.Status != model.PollStatusDraft {
		return nil, errors.New("cannot edit poll that is not in draft status")
	}

	poll.Name = d.Name
	poll.AggregationType = d.AggregationType
	poll.Options = d.Options
	err = repository.UpdatePollItem(poll)
	if err != nil {
		return nil, err
	}

	return viewpoll.ToViewPollModel(poll, nil), nil
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
