package castvote

import "github.com/Samour/voting/site"

type rankedChoiceUpdate struct {
	Select   int
	Unselect int
}

type castVoteModel struct {
	// View controls
	MayVote bool
	// Data
	PollId              string
	PollName            string
	PollAggregationType string
	// Components
	SiteModel     site.SiteModel
	VoteFormModel voteFormModel
}

type voteFormModel struct {
	Voted                 int
	Authenticated         bool
	FptpVoteModel         *fptpVoteModel
	RankedChoiceVoteModel *rankedChoiceVoteModel
}

type fptpVoteModel struct {
	Voted        int
	ErrorMessage string
	PollOptions  []string
}

type rankedChoiceVoteModel struct {
	Voted        int
	PollId       string
	ErrorMessage string
	Rco          rankedChoiceOptions
}

type rankedChoiceOptions struct {
	Selected   []voteOption
	Unselected []voteOption
}

type voteOption struct {
	Index  int
	Option string
}
