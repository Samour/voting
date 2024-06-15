package castvote

import "github.com/Samour/voting/polls/model"

type rankedChoiceUpdate struct {
	Select   int
	Unselect int
}

type castVoteModel struct {
	Poll    *model.Poll
	Rco     *rankedChoiceOptions
	MayVote bool
	Voted   int
}

type rankedChoiceOptions struct {
	Selected   []voteOption
	Unselected []voteOption
}

type voteOption struct {
	Index  int
	Option string
}
