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
	//http.HandleFunc("/about", aboutHandler)

}

//http request handler
type RequestHandler struct {
	Context    *appengine.Context
	User       *user.User
	Google_url string
}

func NewRequestHandler(w http.ResponseWriter, r *http.Request) *RequestHandler {
	handler := new(RequestHandler)

	c := appengine.NewContext(r)

	if c != nil {

		u := user.Current(c)

		if c != nil {
			handler.Context = &c
		}

		if u != nil {
			handler.User = u
			handler.Google_url, _ = user.LogoutURL(c, "/")

		} else {
			handler.Google_url, _ = user.LoginURL(c, "/")

		}

	}

	return handler

}

//handle request
func (handler *RequestHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	c := &handler.Context

	c := c.Infof("HandleRequest")

	var mainTemplate = template.Must(template.ParseFiles(BASE_TEMPLATE, BASE_HEADER_TEMPLATE, BASE_FOOTER_TEMPLATE,
		"./webwallet/templates/main.html"))

	var templateParams = map[string]interface{}{
		"title": "App title",
		"url":   handler.Google_url,
		"user":  handler.User,
	}

	if err := mainTemplate.Execute(w, templateParams); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//main handler
func mainHandler(w http.ResponseWriter, r *http.Request) {
	handler := NewRequestHandler(w, r)

	if &handler.Context != nil {
		handler.HandleRequest(w, r)
	}

}
