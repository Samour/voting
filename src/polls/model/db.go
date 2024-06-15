package model

const DiscriminatorPoll = "poll"
const DiscriminatorVote = "vote:"
const DiscriminatorResult = "result"

type Poll struct {
	PollId          string
	Discriminator   string
	Status          string
	AggregationType string
	Name            string
	Options         []string
	Statistics      VotingStatistics
}

const PollStatusDraft = "draft"
const PollStatusVoting = "voting"
const PollStatusClosed = "closed"

const PollAggregationTypeFirstPastThePost = "fptp"
const PollAggregationTypeRankedChoice = "rankedchoice"

type VotingStatistics struct {
	OpenedAt string
	Votes    int
	ClosedAt string
}

type FptpVote struct {
	PollId        string
	Discriminator string
	Option        int
	CastAt        string
}

type RankedChoiceVote struct {
	PollId        string
	Discriminator string
	Ranked        []int
	CastAt        string
}

type PollResult struct {
	PollId        string
	Discriminator string
	Votes         []OptionVoteCount
}

type OptionVoteCount struct {
	Option    string
	VoteCount int
}
