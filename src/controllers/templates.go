package controllers

import "html/template"

var Templates = template.Must(template.ParseFiles(
	"../resources/components/page_footer.html",
	"../resources/components/page_header.html",

	"../resources/pages/home.html",
	"../resources/pages/new_poll.html",
))
