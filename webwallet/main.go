package webwallet

import (
	"appengine"
	"appengine/user"
	"html/template"
	"net/http"
)

const (
	BASE_TEMPLATE        string = "./webwallet/templates/base.html"
	BASE_HEADER_TEMPLATE string = "./webwallet/templates/baseHeader.html"
	BASE_FOOTER_TEMPLATE string = "./webwallet/templates/baseFooter.html"
)

//initialize function
func init() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/about", aboutHandler)

}

/*
type InterfaceRequestHandler interface {
	HandleRequest(w http.ResponseWriter, r *http.Request)
}
*/

//http request handler
type RequestHandler struct {
	MyContext  appengine.Context
	MyUser     user.User
	Google_url string
}

//handle request
func (handler *RequestHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {

}

//main handler
func mainHandler(w http.ResponseWriter, r *http.Request) {
	var handler RequestHandler = new(RequestHandler)

	c := appengine.NewContext(r)
	handler.MyContext = &c

	u := user.Current(handler.Context)
	handler.MyUser = &u

	if handler.User == nil {
		handler.Google_url, _ = user.LoginURL(handler.Context, "/")

	} else {
		handler.Google_url, _ = user.LogoutURL(handler.Context, "/")

	}

	handler.HandleRequest(w, r)

	var mainTemplate = template.Must(template.ParseFiles(BASE_TEMPLATE, BASE_HEADER_TEMPLATE, BASE_FOOTER_TEMPLATE,
		"./webwallet/templates/main.html"))

	var templateParams = map[string]interface{}{
		"title": "App title",
		"url":   handler.Google_url,
		"user":  handler.User,
	}

	c.Infof("!templateParams[\"url\"]: %s", templateParams["url"])

	if err := mainTemplate.Execute(w, templateParams); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

//about Handler
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	u := user.Current(c)

	var url string

	if u == nil {
		url, _ = user.LoginURL(c, "/")

	} else {
		url, _ = user.LogoutURL(c, "/")

	}

	var mainTemplate = template.Must(template.ParseFiles(BASE_TEMPLATE, BASE_HEADER_TEMPLATE, BASE_FOOTER_TEMPLATE,
		"./webwallet/templates/main.html"))

	var templateParams = map[string]interface{}{
		"title": "App title",
		"url":   url,
		"user":  u,
	}

	c.Infof("!templateParams[\"url\"]: %s", templateParams["url"])

	if err := mainTemplate.Execute(w, templateParams); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
