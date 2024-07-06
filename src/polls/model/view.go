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
	RenderResult        bool
	PollAggregationType string
	Options             []string
	FptpResultModel     ViewPollFptpResultModel
	RcvResultModel      ViewPollRcvResultModel
}

type ViewPollFptpResultModel struct {
	Result []FptpOptionVoteCount
}

type ViewPollRcvResultModel struct {
}

type ViewPollNavigationModel struct {
	PollStatus string
	PollId     string
	VotesCast  int
}
