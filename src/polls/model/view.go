package model

type ViewPollModel struct {
	Poll                 *Poll
	FptpResult           *FptpPollResult
	StatusLabel          string
	AggregationTypeLabel string
	RenderResult         bool
	PollForUpdate        bool
	RenderFullPage       bool
}
