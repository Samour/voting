package polls

type Poll struct {
	PollId        string
	Discriminator string
	Status        string
	Name          string
	Options       []string
}
