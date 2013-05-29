package webwallet

import (
	"appengine"
	"appengine/user"
	"html/template"
	"net/http"
)

//initialize function
func init() {
	http.HandleFunc("/", main)
}

//main handler
func main(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	var url string

	if u == nil {
		url, _ = user.LoginURL(c, "/")

	} else {
		url, _ = user.LogoutURL(c, "/")

	}

	var mainTemplate = template.Must(template.ParseFiles("./webwallet/templates/base.html",
		"./webwallet/templates/main.html"))

	var templateParams = map[string]interface{}{
		"url":  url,
		"user": u,
	}

	c.Infof("!templateParams[\"url\"]: ", templateParams["url"])

	if err := mainTemplate.Execute(w, templateParams); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
}
