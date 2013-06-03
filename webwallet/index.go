package webwallet

import (
	"net/http"
	"webwallet/requesthandler"
)

type IndexHandler struct {
	Config requesthandler.TemplateViewConf
}

func (handler IndexHandler) HandleRequest(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	var params = map[string]interface{}{
		"title": "webcoin",
	}

	return params
}

func (handler IndexHandler) GetConfig() requesthandler.TemplateViewConf {
	return handler.Config
}

//initialize function
func init() {
	var mainHandler IndexHandler
	mainHandler.Config.Url = "/"
	mainHandler.Config.TemplateName = "index"

	requesthandler.MakeRequestHandler(mainHandler)

}
