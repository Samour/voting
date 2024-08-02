package site

type SiteModel struct {
	User User
}

type User struct {
	Authenticated bool
	UserId        string
	DisplayName   string
	UserName      string
}
