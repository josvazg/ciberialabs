package webwallet

import (
	"webwallet/requesthandler"
)

type AboutHandler struct {
	Config requesthandler.TemplateViewConf
}

func (handler AboutHandler) HandleRequest() map[string]interface{} {
	var params = map[string]interface{}{
		"title": "about",
	}

	return params
}

func (handler AboutHandler) GetConfig() requesthandler.TemplateViewConf {
	return handler.Config 
}


//initialize function
func init() {
	var aboutHandler AboutHandler
	aboutHandler.Config.Url = "/about"
	aboutHandler.Config.TemplateName = "about"

	requesthandler.MakeRequestHandler(aboutHandler)

}
