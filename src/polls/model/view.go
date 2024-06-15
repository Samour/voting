package model

type ViewPollModel struct {
	Poll           *Poll
	Result         *PollResult
	StatusLabel    string
	RenderResult   bool
	PollForUpdate  bool
	RenderFullPage bool
}
