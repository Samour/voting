package editpoll

import (
	"github.com/Samour/voting/auth"
	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/site"
)

func buildEditPollModel(s auth.Session, poll model.Poll) editPollModel {
	return editPollModel{
		PollId:              poll.PollId,
		PollName:            poll.Name,
		PollAggregationType: poll.AggregationType,
		MayEdit:             poll.Status == model.PollStatusDraft,
		SiteModel:           site.BuildSiteModel(s),
		OptionsModel:        buildEditPollOptionsModel(poll.Options),
	}
}

func buildEditPollOptionsModel(options []string) editPollOptionsModel {
	return editPollOptionsModel{
		Options: options,
	}
}
