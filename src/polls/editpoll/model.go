package editpoll

import "github.com/Samour/voting/site"

type editPollModel struct {
	PollId              string
	PollName            string
	PollAggregationType string
	MayEdit             bool
	SiteModel           site.SiteModel
	OptionsModel        editPollOptionsModel
}

type editPollOptionsModel struct {
	Options []string
}

type pollDetails struct {
	Name            string
	AggregationType string
	Options         []string
}

type pollOptionsUpdate struct {
	Add    bool
	Remove int
}
