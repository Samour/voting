package controllers

import "html/template"

var templates = template.Must(template.ParseFiles(
	"../resources/components/page_footer.html",
	"../resources/components/page_header.html",

	"../resources/pages/edit_poll.html",
	"../resources/pages/error.html",
	"../resources/pages/home.html",
))
