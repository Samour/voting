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
		ServeViewPoll:          middleware.AuthenticatedWithRedirect(viewpoll.ServeViewPoll),
		HandlePollStatusChange: middleware.AuthenticatedWithError(viewpoll.HandlePollStatusChange),

		ServeNewPoll: middleware.AuthenticatedWithRedirect(createpoll.ServeNewPoll),

		ServeEditPoll:   middleware.AuthenticatedWithRedirect(editpoll.ServeEditPoll),
		ServeSavePoll:   middleware.AuthenticatedWithError(editpoll.ServeSavePoll),
		HandlePatchPoll: middleware.AuthenticatedWithError(editpoll.HandlePatchPoll),

		ServeVotePoll:              castvote.ServeVotePoll,
		HandleCastFptpVote:         castvote.HandleCastFptpVote,
		HandlePatchRankedChoice:    castvote.HandlePatchRankedChoice,
		HandleCastRankedChoiceVote: castvote.HandleCastRankedChoiceVote,
	}
}
