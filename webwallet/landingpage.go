package webwallet

import (
    "appengine"
    "appengine/datastore"
    "appengine/user"
    "html/template"
    "net/http"
    "time"
)

type Greeting struct {
    Author  string
    Content string
    Date    time.Time
}

func init() {
    http.HandleFunc("/", root)
    http.HandleFunc("/sign", sign)
}

func root(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    q := datastore.NewQuery("Greeting").Order("-Date").Limit(10)
    greetings := make([]Greeting, 0, 10)
    if _, err := q.GetAll(c, &greetings); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
	var guestbookTemplate = template.Must(template.ParseFiles("./webwallet/templates/guestbook.html"))

    if err := guestbookTemplate.Execute(w, greetings); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}




func sign(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    g := Greeting{
        Content: r.FormValue("content"),
        Date:    time.Now(),
    }
    if u := user.Current(c); u != nil {
        g.Author = u.String()
    }
    _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Greeting", nil), &g)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/", http.StatusFound)
}