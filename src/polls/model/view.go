package model

type ViewPollModel struct {
	Poll            *Poll
	PollResult      *PollResult
	StatusLabel     string
	OobStatusUpdate bool
}
