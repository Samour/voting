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
	poll := model.Poll{}
	err := repository.GetPollItem(id, model.DiscriminatorPoll, &poll)
	if err != nil {
		return render.HttpResponse{}, err
	}
	if len(poll.PollId) == 0 {
		return render.HttpResponse{
			HttpCode:     http.StatusNotFound,
			ErrorMessage: "Poll not found",
		}, nil
	}

	var fptpResult *model.FptpPollResult
	var rankedChoiceResult *model.RankedChoicePollResult
	if poll.Status == model.PollStatusClosed {
		if poll.AggregationType == model.PollAggregationTypeFirstPastThePost {
			fptpResult, err = loadFptpResult(id)
			if err != nil {
				return render.HttpResponse{}, err
			}
		} else if poll.AggregationType == model.PollAggregationTypeRankedChoice {
			rankedChoiceResult, err = loadRankedChoiceResult(id)
			if err != nil {
				return render.HttpResponse{}, err
			}
		}
	}

	return render.HttpResponse{
		Model: BuildViewPollModel(poll, fptpResult, rankedChoiceResult, renderFullPage),
	}, nil
}

func loadFptpResult(pollId string) (*model.FptpPollResult, error) {
	result := &model.FptpPollResult{}
	err := repository.GetPollItem(pollId, model.DiscriminatorResult, result)
	if err != nil {
		return nil, err
	}
	if len(result.PollId) == 0 {
		return nil, nil
	}

	return result, nil
}

func loadRankedChoiceResult(pollId string) (*model.RankedChoicePollResult, error) {
	result := &model.RankedChoicePollResult{}
	err := repository.GetPollItem(pollId, model.DiscriminatorResult, result)
	if err != nil {
		return nil, err
	}
	if len(result.PollId) == 0 {
		return nil, nil
	}

	return result, nil
}

func updateStatus(id string, status string) (render.HttpResponse, error) {
	poll := model.Poll{}
	err := repository.GetPollItem(id, model.DiscriminatorPoll, &poll)
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
		Model: BuildViewPollModel(poll, nil, nil, false),
	}, nil
}
