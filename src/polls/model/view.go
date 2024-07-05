package model

type ViewPollModel struct {
	// HTML fragment behaviours
	RenderFullPage bool
	PollForUpdate  bool
	// Data
	PollId               string
	PollName             string
	StatusLabel          string
	AggregationTypeLabel string
	// Components
	OptionsModel    ViewPollOptionsModel
	NavigationModel ViewPollNavigationModel
}

type ViewPollOptionsModel struct {
	RenderResult bool
	Result       []FptpOptionVoteCount
	Options      []string
}

type ViewPollNavigationModel struct {
	PollStatus string
	PollId     string
	VotesCast  int
}
