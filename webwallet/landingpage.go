package webwallet

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

/*var templates = template.Must(template.ParseFiles("html/browse.html", "html/add.html",
    "html/templates.html","html/style.css"))*/

func init() {
/*    http.HandleFunc("/login", login)
    http.HandleFunc("/add",add)*/
    http.HandleFunc("/", landingpage)
}

func landingpage(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
}

