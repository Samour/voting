package polls

func FetchPoll(id string) (*Poll, error) {
	return getPollItem(id)
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
