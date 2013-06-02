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
	HandleRequest() map[string]interface{}
	GetConfig() TemplateViewConf
	
}

type TemplateViewConf struct {
	Url          string
	TemplateName string
}


//http request handler
type RequestHandler struct {
	User         *user.User
	Context      appengine.Context
	Google_url   string
	TemplateName string
}


//make request handler
func MakeRequestHandler(templateHandler TemplateViewHandler) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		handler := NewRequestHandler(templateHandler.GetConfig().TemplateName)

		if &handler.Context != nil {
			handler.HandleRequest(w, r)
		}
	}

	http.HandleFunc(templateHandler.GetConfig().Url, handlerFunc)

}

//new request handler
func NewRequestHandler(templateName string) *RequestHandler {
	handler := new(RequestHandler)

	handler.TemplateName = templateName

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
			handler.Google_url, _ = user.LoginURL(c, "/")

		}

	}

	var mainTemplate = template.Must(template.ParseFiles(BASE_TEMPLATE, BASE_HEADER_TEMPLATE, BASE_FOOTER_TEMPLATE,
		"./webwallet/templates/"+handler.TemplateName+".html"))

	var templateParams = map[string]interface{}{
		"title": "App title",
		"url":   handler.Google_url,
		"user":  handler.User,
	}

	if err := mainTemplate.Execute(w, templateParams); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
