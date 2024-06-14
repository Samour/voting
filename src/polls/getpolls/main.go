package getpolls

import (
	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/polls/repository"
)

func FetchAllPolls() ([]model.Poll, error) {
	return repository.ScanPollItems()
}
