package requesthandler

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

//template view handler
type TemplateViewHandler interface {
	HandleRequest(w http.ResponseWriter, r *http.Request) map[string]interface{}
	GetConfig() TemplateViewConf
}

type TemplateViewConf struct {
	Url           string
	TemplateName  string
	HasToBeLogged bool
}

//http request handler
type RequestHandler struct {
	User       *user.User
	Context    appengine.Context
	Google_url string
	Template   TemplateViewHandler
}

//make request handler
func MakeRequestHandler(templateHandler TemplateViewHandler) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		handler := NewRequestHandler(templateHandler)

		if &handler.Context != nil {
			handler.HandleRequest(w, r)
		}
	}

	http.HandleFunc(templateHandler.GetConfig().Url, handlerFunc)

}

//new request handler
func NewRequestHandler(template TemplateViewHandler) *RequestHandler {
	handler := new(RequestHandler)

	handler.Template = template

	return handler

}

//handle request
func (handler *RequestHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	if c != nil {

		handler.Context = c

		u := user.Current(c)

		if u != nil {
			handler.User = u
			handler.Google_url, _ = user.LogoutURL(c, "/")

		} else {
			if handler.Template.GetConfig().HasToBeLogged {
				c.Warningf("user is null: %s", u)

				http.Redirect(w, r, "/", 302)
				return

			} else {
				handler.Google_url, _ = user.LoginURL(c, "/")
			}

		}

	}

	var mainTemplate = template.Must(template.ParseFiles(BASE_TEMPLATE, BASE_HEADER_TEMPLATE, BASE_FOOTER_TEMPLATE,
		"./webwallet/templates/"+handler.Template.GetConfig().TemplateName+".html"))

	var applicationParams = map[string]interface{}{
		"title": "App title",
		"url":   handler.Google_url,
		"user":  handler.User,
	}

	var templateParams = handler.Template.HandleRequest(w, r)

	for k, _ := range templateParams {
		applicationParams[k] = templateParams[k]
	}

	if err := mainTemplate.Execute(w, applicationParams); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
