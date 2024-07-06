package viewpoll

import (
	"errors"
	"net/http"
	"time"

	"github.com/Samour/voting/polls/countvotes"
	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/polls/repository"
	"github.com/Samour/voting/render"
)

func getPoll(id string, renderFullPage bool) (render.HttpResponse, error) {
	poll := &model.Poll{}
	err := repository.GetPollItem(id, model.DiscriminatorPoll, poll)
	if err != nil {
		return render.HttpResponse{}, err
	}
	if len(poll.PollId) == 0 {
		return render.HttpResponse{
			HttpCode:     http.StatusNotFound,
			ErrorMessage: "Poll not found",
		}, nil
	}

	pollResult := &model.FptpPollResult{}
	err = repository.GetPollItem(id, model.DiscriminatorResult, pollResult)
	if err != nil {
		return render.HttpResponse{}, err
	}
	if len(pollResult.PollId) == 0 {
		pollResult = nil
	}

	return render.HttpResponse{
		Model: BuildViewPollModel(poll, pollResult, renderFullPage),
	}, nil
}

func updateStatus(id string, status string) (render.HttpResponse, error) {
	poll := &model.Poll{}
	err := repository.GetPollItem(id, model.DiscriminatorPoll, poll)
	if err != nil {
		return render.HttpResponse{}, err
	}
	if len(poll.PollId) == 0 {
		return render.HttpResponse{
			HttpCode:     http.StatusNotFound,
			ErrorMessage: "Poll not found",
		}, nil
	}

	startVoteCounting := false
	if status == model.PollStatusVoting {
		if poll.Status != model.PollStatusDraft {
			return render.HttpResponse{
				HttpCode:     http.StatusBadRequest,
				ErrorMessage: "Cannot open voting on poll",
			}, nil
		}
		poll.Statistics.OpenedAt = time.Now().In(time.UTC).Format(time.RFC3339)
	} else if status == model.PollStatusClosed {
		if poll.Status != model.PollStatusVoting {
			return render.HttpResponse{
				HttpCode:     http.StatusBadRequest,
				ErrorMessage: "Voting is not currently open on poll",
			}, nil
		}
		poll.Statistics.ClosedAt = time.Now().In(time.UTC).Format(time.RFC3339)
		startVoteCounting = true
	} else {
		return render.HttpResponse{}, errors.New("unknown status")
	}

	poll.Status = status
	err = repository.UpdatePollItem(poll)
	if err != nil {
		return render.HttpResponse{}, err
	}

	if startVoteCounting {
		go countvotes.CountVotes(id)
	}

	return render.HttpResponse{
		Model: BuildViewPollModel(poll, nil, false),
	}, nil
}
