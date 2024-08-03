package editpoll

import (
	"net/http"

	"github.com/Samour/voting/auth"
	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/polls/repository"
	"github.com/Samour/voting/polls/viewpoll"
	"github.com/Samour/voting/render"
)

func getPoll(id string) (render.HttpResponse, error) {
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

	return render.HttpResponse{
		Model: buildEditPollModel(poll),
	}, nil
}

func updatePollDetails(s auth.Session, id string, d pollDetails) (render.HttpResponse, error) {
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

	if poll.Status != model.PollStatusDraft {
		return render.HttpResponse{
			HttpCode:     http.StatusBadRequest,
			ErrorMessage: "Cannot edit poll that is not in draft status",
		}, nil
	}

	poll.Name = d.Name
	poll.AggregationType = d.AggregationType
	poll.Options = d.Options
	err = repository.UpdatePollItem(poll)
	if err != nil {
		return render.HttpResponse{}, err
	}

	return render.HttpResponse{
		Model: viewpoll.BuildViewPollModel(s, viewpoll.ViewPollData{
			Poll: poll,
		}),
	}, nil
}

func patchPollOptions(options []string, u pollOptionsUpdate) render.HttpResponse {
	if u.Remove >= 0 && u.Remove < len(options) {
		options = append(options[:u.Remove], options[u.Remove+1:]...)
	}
	if u.Add || len(options) == 0 {
		options = append(options, "")
	}

	return render.HttpResponse{
		Model: buildEditPollOptionsModel(options),
	}
}
