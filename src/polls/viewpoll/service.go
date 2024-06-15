package viewpoll

import (
	"errors"
	"time"

	"github.com/Samour/voting/polls/countvotes"
	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/polls/repository"
)

func getPoll(id string) (*model.ViewPollModel, error) {
	poll := &model.Poll{}
	err := repository.GetPollItem(id, model.DiscriminatorPoll, poll)
	if err != nil {
		return nil, err
	}
	if len(poll.PollId) == 0 {
		return nil, nil
	}

	pollResult := &model.FptpPollResult{}
	err = repository.GetPollItem(id, model.DiscriminatorResult, pollResult)
	if err != nil {
		return nil, err
	}
	if len(pollResult.PollId) == 0 {
		pollResult = nil
	}

	return ToViewPollModel(poll, pollResult), nil
}

func ToViewPollModel(p *model.Poll, r *model.FptpPollResult) *model.ViewPollModel {
	statusLabel := p.Status
	pollForUpdate := false
	if p.Status == model.PollStatusClosed && r == nil {
		statusLabel = "closed; vote count pending"
		pollForUpdate = true
	}

	aggregationTypeLabel := ""
	if p.AggregationType == model.PollAggregationTypeFirstPastThePost {
		aggregationTypeLabel = "First past the post"
	} else if p.AggregationType == model.PollAggregationTypeRankedChoice {
		aggregationTypeLabel = "Ranked choice"
	}

	return &model.ViewPollModel{
		Poll:                 p,
		FptpResult:           r,
		StatusLabel:          statusLabel,
		AggregationTypeLabel: aggregationTypeLabel,
		RenderResult:         r != nil,
		PollForUpdate:        pollForUpdate,
		RenderFullPage:       true,
	}
}

func updateStatus(id string, status string) (*model.ViewPollModel, error) {
	poll := &model.Poll{}
	err := repository.GetPollItem(id, model.DiscriminatorPoll, poll)
	if err != nil {
		return nil, err
	}
	if len(poll.PollId) == 0 {
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

	return ToViewPollModel(poll, nil), nil
}
