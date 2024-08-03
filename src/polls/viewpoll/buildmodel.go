package viewpoll

import (
	"fmt"

	"github.com/Samour/voting/auth"
	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/site"
)

type ViewPollData struct {
	Poll           model.Poll
	FptpResult     *model.FptpPollResult
	RcvResult      *model.RankedChoicePollResult
	RenderFullPage bool
}

func BuildViewPollModel(s auth.Session, d ViewPollData) model.ViewPollModel {
	statusLabel := d.Poll.Status
	pollForUpdate := false
	if d.Poll.Status == model.PollStatusClosed && d.FptpResult == nil && d.RcvResult == nil {
		statusLabel = "closed; vote count pending"
		pollForUpdate = true
	}

	aggregationTypeLabel := ""
	if d.Poll.AggregationType == model.PollAggregationTypeFirstPastThePost {
		aggregationTypeLabel = "First past the post"
	} else if d.Poll.AggregationType == model.PollAggregationTypeRankedChoice {
		aggregationTypeLabel = "Ranked choice"
	}

	return model.ViewPollModel{
		RenderFullPage:       d.RenderFullPage,
		PollForUpdate:        pollForUpdate,
		PollId:               d.Poll.PollId,
		PollName:             d.Poll.Name,
		StatusLabel:          statusLabel,
		AggregationTypeLabel: aggregationTypeLabel,
		SiteModel:            site.BuildSiteModel(s),
		OptionsModel:         buildViewPollOptionsModel(d.Poll, d.FptpResult, d.RcvResult),
		NavigationModel:      buildViewPollNavigationModel(d.Poll),
	}
}

func buildViewPollOptionsModel(p model.Poll, fptp *model.FptpPollResult, rcv *model.RankedChoicePollResult) model.ViewPollOptionsModel {
	return model.ViewPollOptionsModel{
		RenderResult:        fptp != nil || rcv != nil,
		PollAggregationType: p.AggregationType,
		Options:             p.Options,
		FptpResultModel:     buildViewPollFptpResultModel(fptp),
		RcvResultModel:      buildViewPollRcvResultModel(rcv),
	}
}

func buildViewPollFptpResultModel(r *model.FptpPollResult) model.ViewPollFptpResultModel {
	var result []model.FptpOptionVoteCount
	if r != nil {
		result = r.Votes
	}

	return model.ViewPollFptpResultModel{
		Result: result,
	}
}

func buildViewPollRcvResultModel(r *model.RankedChoicePollResult) model.ViewPollRcvResultModel {
	if r == nil {
		return model.ViewPollRcvResultModel{}
	}

	roundCount := 1
	for _, v := range r.Votes {
		if len(v.RoundVotes) > roundCount {
			roundCount = len(v.RoundVotes)
		}
	}
	roundTitles := make([]string, roundCount)
	for i := 0; i < roundCount; i++ {
		roundTitles[i] = fmt.Sprint(i + 1)
	}

	votes := make([]model.RankedChoiceOptionVoteCount, len(r.Votes))
	for i, v := range r.Votes {
		c := model.RankedChoiceOptionVoteCount{
			Option:     v.Option,
			RoundVotes: make([]int, roundCount),
		}
		for j := 0; j < roundCount; j++ {
			if j < len(v.RoundVotes) {
				c.RoundVotes[j] = v.RoundVotes[j]
			} else {
				c.RoundVotes[j] = -1
			}
		}
		votes[i] = c
	}

	return model.ViewPollRcvResultModel{
		RoundTitles: roundTitles,
		Result:      votes,
	}
}

func buildViewPollNavigationModel(p model.Poll) model.ViewPollNavigationModel {
	return model.ViewPollNavigationModel{
		PollStatus: p.Status,
		PollId:     p.PollId,
		VotesCast:  p.Statistics.Votes,
	}
}
