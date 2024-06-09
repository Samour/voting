package polls

type PollDetails struct {
	Name    string
	Options []string
}

type PollOptionsUpdate struct {
	Details PollDetails
	Add     bool
	Remove  int
}

func UpdatePollDetails(id string, name string, options []string) (*Poll, error) {
	poll, err := getPollItem(id)
	if err != nil {
		return nil, err
	}
	if poll == nil {
		return nil, nil
	}

	// TODO handle poll that is not in draft

	poll.Name = name
	poll.Options = options
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

	return UpdatePollDetails(id, d.Name, d.Options)
}
