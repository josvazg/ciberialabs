package webwallet

import (
	"net/http"
	"webwallet/requesthandler"
)

type ContactHandler struct {
	Config requesthandler.TemplateViewConf
}

func (handler ContactHandler) HandleRequest(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	var params = map[string]interface{}{
		"title": "contact",
	}

	return params
}

func (handler ContactHandler) GetConfig() requesthandler.TemplateViewConf {
	return handler.Config
}

//initialize function
func init() {
	var aboutHandler ContactHandler
	aboutHandler.Config.Url = "/contact"
	aboutHandler.Config.TemplateName = "contact"

	requesthandler.MakeRequestHandler(aboutHandler)

}
