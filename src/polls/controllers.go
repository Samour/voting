package polls

import (
	"github.com/Samour/voting/middleware"
	"github.com/Samour/voting/polls/castvote"
	"github.com/Samour/voting/polls/createpoll"
	"github.com/Samour/voting/polls/editpoll"
	"github.com/Samour/voting/polls/viewpoll"
	"github.com/Samour/voting/types"
)

type PollControllers struct {
	ServeViewPoll          types.Controller
	HandlePollStatusChange types.Controller

	ServeNewPoll types.Controller

	ServeEditPoll   types.Controller
	ServeSavePoll   types.Controller
	HandlePatchPoll types.Controller

	ServeVotePoll              types.Controller
	HandleCastFptpVote         types.Controller
	HandlePatchRankedChoice    types.Controller
	HandleCastRankedChoiceVote types.Controller
}

func CreatePollControllers() PollControllers {
	return PollControllers{
		ServeViewPoll:          middleware.RedirectUnauthenticated(viewpoll.ServeViewPoll),
		HandlePollStatusChange: middleware.PreventUnauthenticated(viewpoll.HandlePollStatusChange),

		ServeNewPoll: middleware.RedirectUnauthenticated(createpoll.ServeNewPoll),

		ServeEditPoll:   middleware.RedirectUnauthenticated(editpoll.ServeEditPoll),
		ServeSavePoll:   middleware.PreventUnauthenticated(editpoll.ServeSavePoll),
		HandlePatchPoll: middleware.PreventUnauthenticated(editpoll.HandlePatchPoll),

		ServeVotePoll:              castvote.ServeVotePoll,
		HandleCastFptpVote:         castvote.HandleCastFptpVote,
		HandlePatchRankedChoice:    castvote.HandlePatchRankedChoice,
		HandleCastRankedChoiceVote: castvote.HandleCastRankedChoiceVote,
	}
}
