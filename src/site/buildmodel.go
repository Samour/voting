package site

import "github.com/Samour/voting/auth"

func BuildSiteModel(session auth.Session) SiteModel {
	user := User{}
	if len(session.User.UserId) > 0 {
		user = User{
			Authenticated: true,
			UserId:        session.User.UserId,
			DisplayName:   session.User.DisplayName,
			UserName:      session.User.Username,
		}
	}

	return SiteModel{
		User: user,
	}
}
