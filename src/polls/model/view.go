package model

type ViewPollModel struct {
	Poll            *Poll
	Result          *PollResult
	StatusLabel     string
	RenderResult    bool
	OobStatusUpdate bool
}
