package viewpoll

import (
	"fmt"

	"github.com/Samour/voting/polls/model"
)

func BuildViewPollModel(p model.Poll, fptp *model.FptpPollResult, rcv *model.RankedChoicePollResult, renderFullPage bool) model.ViewPollModel {
	statusLabel := p.Status
	pollForUpdate := false
	if p.Status == model.PollStatusClosed && fptp == nil && rcv == nil {
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
		OptionsModel:         buildViewPollOptionsModel(p, fptp, rcv),
		NavigationModel:      buildViewPollNavigationModel(p),
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
