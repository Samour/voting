package polls

type PollItem struct {
	PollId        string
	Discriminator string
	Status        string
	Name          string
	Options       []string
}
