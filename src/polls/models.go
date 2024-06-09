package polls

type Poll struct {
	PollId        string
	Discriminator string
	Status        string
	Name          string
	Options       []string
	Statistics    VotingStatistics
}

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
