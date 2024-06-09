package polls

func FetchPoll(id string) (*Poll, error) {
	return getPollItem(id)
}

func FetchAllPolls() ([]Poll, error) {
	return scanPollItems()
}
