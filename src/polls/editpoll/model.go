package editpoll

type editPollModel struct {
	PollId              string
	PollName            string
	PollAggregationType string
	MayEdit             bool
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
