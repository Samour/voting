package polls

func FetchPoll(id string) (*Poll, error) {
	return getPollItem(id)
}
