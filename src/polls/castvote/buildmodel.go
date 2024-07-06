package castvote

import "github.com/Samour/voting/polls/model"

type castVoteData struct {
	Poll         model.Poll
	Voted        int
	Ranked       []int
	ErrorMessage string
}

func buildCastVoteModel(d castVoteData) castVoteModel {
	return castVoteModel{
		MayVote:             d.Poll.Status == model.PollStatusVoting,
		PollId:              d.Poll.PollId,
		PollName:            d.Poll.Name,
		PollAggregationType: d.Poll.AggregationType,
		VoteFormModel:       buildVoteFormModel(d),
	}
}

func buildVoteFormModel(d castVoteData) voteFormModel {
	var fptpModel *fptpVoteModel
	var rcvModel *rankedChoiceVoteModel

	if d.Poll.AggregationType == model.PollAggregationTypeFirstPastThePost {
		m := buildFptpVoteModel(d)
		fptpModel = &m
	} else if d.Poll.AggregationType == model.PollAggregationTypeRankedChoice {
		m := buildRankedChoiceVoteModel(d)
		rcvModel = &m
	}

	return voteFormModel{
		Voted:                 d.Voted,
		FptpVoteModel:         fptpModel,
		RankedChoiceVoteModel: rcvModel,
	}
}

func buildFptpVoteModel(d castVoteData) fptpVoteModel {
	return fptpVoteModel{
		Voted:        d.Voted,
		PollOptions:  d.Poll.Options,
		ErrorMessage: d.ErrorMessage,
	}
}

func buildRankedChoiceVoteModel(d castVoteData) rankedChoiceVoteModel {
	selectedOptions := make([]voteOption, len(d.Ranked))
	for i, o := range d.Ranked {
		selectedOptions[i] = voteOption{
			Index:  o,
			Option: d.Poll.Options[o],
		}
	}

	unselected := []voteOption{}
	for i, o := range d.Poll.Options {
		if !contains(i, d.Ranked) {
			unselected = append(unselected, voteOption{
				Index:  i,
				Option: o,
			})
		}
	}

	return rankedChoiceVoteModel{
		Voted:  d.Voted,
		PollId: d.Poll.PollId,
		Rco: rankedChoiceOptions{
			Selected:   selectedOptions,
			Unselected: unselected,
		},
		ErrorMessage: d.ErrorMessage,
	}
}

func contains(find int, values []int) bool {
	for _, i := range values {
		if i == find {
			return true
		}
	}

	return false
}
