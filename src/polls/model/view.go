package model

import "github.com/Samour/voting/site"

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
	SiteModel       site.SiteModel
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
	RoundTitles []string
	Result      []RankedChoiceOptionVoteCount
}

type ViewPollNavigationModel struct {
	PollStatus string
	PollId     string
	VotesCast  int
}
