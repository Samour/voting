package polls

func FetchPoll(id string) (*PollItem, error) {
	return getPollItem(id)
}
