package castvote

import "github.com/Samour/voting/polls/model"

type castVoteModel struct {
	Poll    *model.Poll
	MayVote bool
	Voted   int
}
