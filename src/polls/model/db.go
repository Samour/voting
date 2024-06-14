package model

const DiscriminatorPoll = "poll"
const DiscriminatorVote = "vote:"

type Poll struct {
	PollId        string
	Discriminator string
	Status        string
	Name          string
	Options       []string
	Statistics    VotingStatistics
}

const PollStatusDraft = "draft"
const PollStatusVoting = "voting"
const PollStatusClosed = "closed"

type VotingStatistics struct {
	OpenedAt string
	Votes    int
	ClosedAt string
}

type Vote struct {
	PollId        string
	Discriminator string
	Option        int
	CastAt        string
}
