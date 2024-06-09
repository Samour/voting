package polls

import "errors"

type PollDetails struct {
	Name    string
	Options []string
}

type PollOptionsUpdate struct {
	Details PollDetails
	Add     bool
	Remove  int
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

func PatchPollOptions(id string, u PollOptionsUpdate) (*Poll, error) {
	d := u.Details
	if u.Remove >= 0 && u.Remove < len(d.Options) {
		d.Options = append(d.Options[:u.Remove], d.Options[u.Remove+1:]...)
	}
	if u.Add || len(d.Options) == 0 {
		d.Options = append(d.Options, "")
	}

	return UpdatePollDetails(id, d)
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
