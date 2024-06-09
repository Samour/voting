package polls

import "github.com/Samour/voting/utils"

func CreatePoll() string {
	id := utils.IdGen()
	return id
}
