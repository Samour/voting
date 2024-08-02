package home

import (
	"github.com/Samour/voting/polls/model"
	"github.com/Samour/voting/site"
)

type homeModel struct {
	SiteModel site.SiteModel
	Polls     []model.Poll
}
