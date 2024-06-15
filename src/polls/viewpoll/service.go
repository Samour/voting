package viewpoll

import (
	"errors"
	"time"

	"github.com/Samour/voting/polls/countvotes"
	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/polls/repository"
)

func getPoll(id string) (*model.ViewPollModel, error) {
	poll, err := repository.GetPollItem(id)
	if err != nil {
		return nil, err
	}
	if poll == nil {
		return nil, nil
	}

	pollResult, err := repository.GetPollResultItem(id)
	if err != nil {
		return nil, err
	}

	return ToViewPollModel(poll, pollResult), nil
}

func ToViewPollModel(p *model.Poll, r *model.PollResult) *model.ViewPollModel {
	statusLabel := p.Status
	if p.Status == model.PollStatusClosed && r == nil {
		statusLabel = "closed; vote count pending"
	}

	return &model.ViewPollModel{
		Poll:            p,
		PollResult:      r,
		StatusLabel:     statusLabel,
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

	startVoteCounting := false
	if status == model.PollStatusVoting {
		if poll.Status != model.PollStatusDraft {
			return nil, errors.New("cannot open voting on poll")
		}
		poll.Statistics.OpenedAt = time.Now().In(time.UTC).Format(time.RFC3339)
	} else if status == model.PollStatusClosed {
		if poll.Status != model.PollStatusVoting {
			return nil, errors.New("voting is not currently open on poll")
		}
		poll.Statistics.ClosedAt = time.Now().In(time.UTC).Format(time.RFC3339)
		startVoteCounting = true
	} else {
		return nil, errors.New("unknown status")
	}

	poll.Status = status
	err = repository.UpdatePollItem(poll)
	if err != nil {
		return nil, err
	}

	if startVoteCounting {
		go countvotes.CountVotes(id)
	}

	model := ToViewPollModel(poll, nil)
	model.OobStatusUpdate = true
	return model, nil
}
