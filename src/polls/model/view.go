package model

type ViewPollModel struct {
	Poll                 *Poll
	Result               *PollResult
	StatusLabel          string
	AggregationTypeLabel string
	RenderResult         bool
	PollForUpdate        bool
	RenderFullPage       bool
}
