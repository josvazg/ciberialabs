package webwallet

import (
	"net/http"
	"webwallet/requesthandler"
)

type MainHandler struct {
	Config requesthandler.TemplateViewConf
}

func (handler MainHandler) HandleRequest(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	var params = map[string]interface{}{
		"title": "webcoin",
	}
	return params
}

func (handler MainHandler) GetConfig() requesthandler.TemplateViewConf {
	return handler.Config
}

//initialize function
func init() {
	var mainHandler MainHandler
	mainHandler.Config.Url = "/main"
	mainHandler.Config.TemplateName = "main"
	mainHandler.Config.HasToBeLogged = true

	requesthandler.MakeRequestHandler(mainHandler)

}
