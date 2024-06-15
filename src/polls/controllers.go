package polls

import (
	"net/http"

	"github.com/Samour/voting/polls/castvote"
	"github.com/Samour/voting/polls/createpoll"
	"github.com/Samour/voting/polls/editpoll"
	"github.com/Samour/voting/polls/viewpoll"
)

type PollControllers struct {
	ServeViewPoll          func(http.ResponseWriter, *http.Request)
	HandlePollStatusChange func(http.ResponseWriter, *http.Request)

	ServeNewPoll func(http.ResponseWriter, *http.Request)

	ServeEditPoll   func(http.ResponseWriter, *http.Request)
	ServeSavePoll   func(http.ResponseWriter, *http.Request)
	HandlePatchPoll func(http.ResponseWriter, *http.Request)

	ServeVotePoll           func(http.ResponseWriter, *http.Request)
	ServeCastFptpVote       func(http.ResponseWriter, *http.Request)
	HandlePatchRankedChoice func(http.ResponseWriter, *http.Request)
}

func CreatePollControllers() PollControllers {
	return PollControllers{
		ServeViewPoll:          viewpoll.ServeViewPoll,
		HandlePollStatusChange: viewpoll.HandlePollStatusChange,

		ServeNewPoll: createpoll.ServeNewPoll,

		ServeEditPoll:   editpoll.ServeEditPoll,
		ServeSavePoll:   editpoll.ServeSavePoll,
		HandlePatchPoll: editpoll.HandlePatchPoll,

		ServeVotePoll:           castvote.ServeVotePoll,
		ServeCastFptpVote:       castvote.ServeCastFptpVote,
		HandlePatchRankedChoice: castvote.HandlePatchRankedChoice,
	}
}
