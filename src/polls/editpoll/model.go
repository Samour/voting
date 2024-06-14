package editpoll

import "github.com/Samour/voting/polls/model"

type editPollModel struct {
	Poll    *model.Poll
	MayEdit bool
}

type pollDetails struct {
	Name    string
	Options []string
}

type pollOptionsUpdate struct {
	Add    bool
	Remove int
}
