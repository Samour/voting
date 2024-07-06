package viewpoll

import "github.com/Samour/voting/polls/model"

func BuildViewPollModel(p *model.Poll, r *model.FptpPollResult, renderFullPage bool) model.ViewPollModel {
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

	return model.ViewPollModel{
		RenderFullPage:       renderFullPage,
		PollForUpdate:        pollForUpdate,
		PollId:               p.PollId,
		PollName:             p.Name,
		StatusLabel:          statusLabel,
		AggregationTypeLabel: aggregationTypeLabel,
		OptionsModel:         buildViewPollOptionsModel(p, r),
		NavigationModel:      buildViewPollNavigationModel(p),
	}
}

func buildViewPollOptionsModel(p *model.Poll, r *model.FptpPollResult) model.ViewPollOptionsModel {
	var result []model.FptpOptionVoteCount = nil
	if r != nil {
		result = r.Votes
	}

	return model.ViewPollOptionsModel{
		RenderResult: r != nil,
		Result:       result,
		Options:      p.Options,
	}
}

func buildViewPollNavigationModel(p *model.Poll) model.ViewPollNavigationModel {
	return model.ViewPollNavigationModel{
		PollStatus: p.Status,
		PollId:     p.PollId,
		VotesCast:  p.Statistics.Votes,
	}
}
