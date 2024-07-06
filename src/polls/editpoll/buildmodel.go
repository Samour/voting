package editpoll

import "github.com/Samour/voting/polls/model"

func buildEditPollModel(poll model.Poll) editPollModel {
	return editPollModel{
		PollId:              poll.PollId,
		PollName:            poll.Name,
		PollAggregationType: poll.AggregationType,
		MayEdit:             poll.Status == model.PollStatusDraft,
		OptionsModel:        buildEditPollOptionsModel(poll.Options),
	}
}

func buildEditPollOptionsModel(options []string) editPollOptionsModel {
	return editPollOptionsModel{
		Options: options,
	}
}
