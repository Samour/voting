package polls

import (
	"errors"
	"time"
)

type PollDetails struct {
	Name    string
	Options []string
}

type PollOptionsUpdate struct {
	Add    bool
	Remove int
}

func UpdatePollDetails(id string, d PollDetails) (*Poll, error) {
	poll, err := getPollItem(id)
	if err != nil {
		return nil, err
	}
	if poll == nil {
		return nil, nil
	}

	// TODO handle poll that is not in draft

	poll.Name = d.Name
	poll.Options = d.Options
	err = updatePollItem(poll)
	if err != nil {
		return nil, err
	}

	return poll, nil
}

func PatchPollOptions(options []string, u PollOptionsUpdate) []string {
	if u.Remove >= 0 && u.Remove < len(options) {
		options = append(options[:u.Remove], options[u.Remove+1:]...)
	}
	if u.Add || len(options) == 0 {
		options = append(options, "")
	}

	return options
}

func UpdateStatus(id string, status string) (*Poll, error) {
	poll, err := getPollItem(id)
	if err != nil {
		return nil, err
	}
	if poll == nil {
		return nil, nil
	}

	if status == "voting" {
		if poll.Status != "draft" {
			return nil, errors.New("cannot open voting on poll")
		}
		poll.Statistics.OpenedAt = time.Now().In(time.UTC).Format(time.RFC3339)
	} else {
		return nil, errors.New("unknown status")
	}

	poll.Status = status
	err = updatePollItem(poll)
	if err != nil {
		return nil, err
	}

	return poll, nil
}
