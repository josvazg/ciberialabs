package webwallet

import (
	"webwallet/requesthandler"
)

type MainHandler struct {
	Config requesthandler.TemplateViewConf
}


func (handler MainHandler) HandleRequest() map[string]interface{} {
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
	mainHandler.Config.Url = "/"
	mainHandler.Config.TemplateName = "main"
	
	requesthandler.MakeRequestHandler(mainHandler)

}
