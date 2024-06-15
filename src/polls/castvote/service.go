package castvote

import (
	"errors"
	"time"

	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/polls/repository"
	"github.com/Samour/voting/utils"
)

func getPollVoteForm(pollId string) (*castVoteModel, error) {
	poll, err := repository.GetPollItem(pollId)
	if err != nil {
		return nil, err
	}
	if poll == nil {
		return nil, nil
	}

	var rco *rankedChoiceOptions = nil
	if poll.AggregationType == model.PollAggregationTypeRankedChoice {
		uo := make([]voteOption, len(poll.Options))
		for i, o := range poll.Options {
			uo[i] = voteOption{
				Index:  i,
				Option: o,
			}
		}
		rco = &rankedChoiceOptions{
			Unselected: uo,
		}
	}

	return &castVoteModel{
		Poll:    poll,
		Rco:     rco,
		MayVote: poll.Status == model.PollStatusVoting,
		Voted:   -1,
	}, nil
}

func castFptpVote(pollId string, option int) (*castVoteModel, error) {
	poll, err := repository.GetPollItem(pollId)
	if err != nil {
		return nil, err
	}
	if poll == nil {
		return nil, errors.New("poll not found")
	}

	if poll.Status != model.PollStatusVoting {
		return nil, errors.New("poll is not open for voting")
	}

	if option < 0 || option >= len(poll.Options) {
		errorMsg := "You must select an option to vote for"
		return &castVoteModel{
			Poll:        poll,
			MayVote:     true,
			Voted:       -1,
			ErrorString: &errorMsg,
		}, nil
	}

	voteId := utils.IdGen()
	discriminator := model.DiscriminatorVote + voteId
	vote := model.Vote{
		PollId:        pollId,
		Discriminator: discriminator,
		Option:        option,
		CastAt:        time.Now().In(time.UTC).Format(time.RFC3339),
	}

	err = repository.RecordVote(&vote)
	if err != nil {
		return nil, err
	}

	return &castVoteModel{
		Poll:    poll,
		MayVote: true,
		Voted:   option,
	}, nil
}

func updateRankedChoiceOption(pollId string, options []int, u rankedChoiceUpdate) (*castVoteModel, error) {
	poll, err := repository.GetPollItem(pollId)
	if err != nil {
		return nil, err
	}

	if u.Unselect >= 0 {
		options = removeFromList(options, u.Unselect)
	}
	if u.Select >= 0 {
		options = append(options, u.Select)
	}

	selected, err := constructVoteOptionsList(poll, options)
	if err != nil {
		return nil, err
	}
	unselected := constructVoteOptionsListWithout(poll, options)

	return &castVoteModel{
		Poll: poll,
		Rco: &rankedChoiceOptions{
			Selected:   selected,
			Unselected: unselected,
		},
		MayVote: true,
		Voted:   -1,
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

func constructVoteOptionsList(poll *model.Poll, options []int) ([]voteOption, error) {
	list := make([]voteOption, len(options))
	for i, o := range options {
		if o < 0 || o >= len(poll.Options) {
			return nil, errors.New("option value out of range")
		}

		list[i] = voteOption{
			Index:  o,
			Option: poll.Options[o],
		}
	}

	return list, nil
}

func constructVoteOptionsListWithout(poll *model.Poll, options []int) []voteOption {
	list := make([]voteOption, 0)
OPTIONS:
	for i, o := range poll.Options {
		for _, ex := range options {
			if i == ex {
				continue OPTIONS
			}
		}

		list = append(list, voteOption{
			Index:  i,
			Option: o,
		})
	}

	return list
}
