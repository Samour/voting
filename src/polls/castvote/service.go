package castvote

import (
	"net/http"
	"time"

	"github.com/Samour/voting/auth"
	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/polls/repository"
	"github.com/Samour/voting/render"
	"github.com/Samour/voting/utils"
)

func getPollVoteForm(s auth.Session, pollId string) (render.HttpResponse, error) {
	poll := model.Poll{}
	err := repository.GetPollItem(pollId, model.DiscriminatorPoll, &poll)
	if err != nil {
		return render.HttpResponse{}, err
	}
	if len(poll.PollId) == 0 {
		return render.HttpResponse{
			HttpCode:     http.StatusNotFound,
			ErrorMessage: "Poll not found",
		}, nil
	}

	return render.HttpResponse{
		Model: buildCastVoteModel(s, castVoteData{
			Poll:  poll,
			Voted: -1,
		}),
	}, nil
}

func castFptpVote(pollId string, option int) (render.HttpResponse, error) {
	poll := model.Poll{}
	err := repository.GetPollItem(pollId, model.DiscriminatorPoll, &poll)
	if err != nil {
		return render.HttpResponse{}, err
	}
	if len(poll.PollId) == 0 {
		return render.HttpResponse{
			HttpCode:     http.StatusNotFound,
			ErrorMessage: "Poll not found",
		}, nil
	}

	if poll.Status != model.PollStatusVoting {
		return render.HttpResponse{
			HttpCode:     http.StatusBadRequest,
			ErrorMessage: "Poll is not open for voting",
		}, nil
	}
	if poll.AggregationType != model.PollAggregationTypeFirstPastThePost {
		return render.HttpResponse{
			HttpCode:     http.StatusBadRequest,
			ErrorMessage: "Incorrect vote type for poll",
		}, nil
	}

	if option < 0 || option >= len(poll.Options) {
		return render.HttpResponse{
			Model: buildVoteFormModel(castVoteData{
				Poll:         poll,
				Voted:        -1,
				ErrorMessage: "You must select an option to vote for",
			}),
		}, nil
	}

	voteId := utils.IdGen()
	discriminator := model.DiscriminatorVote + voteId
	vote := model.FptpVote{
		PollId:        pollId,
		Discriminator: discriminator,
		Option:        option,
		CastAt:        time.Now().In(time.UTC).Format(time.RFC3339),
	}

	err = repository.RecordVote(pollId, &vote)
	if err != nil {
		return render.HttpResponse{}, err
	}

	return render.HttpResponse{
		Model: buildVoteFormModel(castVoteData{
			Poll:  poll,
			Voted: option,
		}),
	}, nil
}

func updateRankedChoiceOption(pollId string, options []int, u rankedChoiceUpdate) (render.HttpResponse, error) {
	poll := model.Poll{}
	err := repository.GetPollItem(pollId, model.DiscriminatorPoll, &poll)
	if err != nil {
		return render.HttpResponse{}, err
	}

	if u.Unselect >= 0 {
		options = removeFromList(options, u.Unselect)
	}
	if u.Select >= 0 {
		options = append(options, u.Select)
	}

	return render.HttpResponse{
		Model: buildRankedChoiceVoteModel(castVoteData{
			Poll:   poll,
			Voted:  -1,
			Ranked: options,
		}),
	}, nil
}

func removeFromList(l []int, v int) []int {
	r := make([]int, 0)
	for _, x := range l {
		if x != v {
			r = append(r, x)
		}
	}

	return r
}

func castRankedChoiceVote(pollId string, ranked []int) (render.HttpResponse, error) {
	poll := model.Poll{}
	err := repository.GetPollItem(pollId, model.DiscriminatorPoll, &poll)
	if err != nil {
		return render.HttpResponse{}, err
	}
	if len(poll.PollId) == 0 {
		return render.HttpResponse{
			HttpCode:     http.StatusNotFound,
			ErrorMessage: "Poll does not exist",
		}, nil
	}

	if poll.Status != model.PollStatusVoting {
		return render.HttpResponse{
			HttpCode:     http.StatusBadRequest,
			ErrorMessage: "Poll is not open for voting",
		}, nil
	}
	if poll.AggregationType != model.PollAggregationTypeRankedChoice {
		return render.HttpResponse{
			HttpCode:     http.StatusBadRequest,
			ErrorMessage: "Incorrect vote type for poll",
		}, nil
	}

	if len(ranked) < len(poll.Options) {
		return render.HttpResponse{
			Model: buildVoteFormModel(castVoteData{
				Poll:         poll,
				Voted:        -1,
				Ranked:       ranked,
				ErrorMessage: "All options must be selected",
			}),
		}, nil
	}

	voteId := utils.IdGen()
	vote := model.RankedChoiceVote{
		PollId:        pollId,
		Discriminator: model.DiscriminatorVote + voteId,
		Ranked:        ranked,
		CastAt:        time.Now().In(time.UTC).Format(time.RFC3339),
	}

	err = repository.RecordVote(pollId, &vote)
	if err != nil {
		return render.HttpResponse{}, err
	}

	return render.HttpResponse{
		Model: buildVoteFormModel(castVoteData{
			Poll:   poll,
			Voted:  1,
			Ranked: ranked,
		}),
	}, nil
}
