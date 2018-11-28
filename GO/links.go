
package main

import (
	"html/template"
	"net/http"
	//"io/ioutil"
	//"fmt"
)

type Link struct {
	URL string
	Name string
}

type LinkData struct {
	Links []Link
}

func LoadLinks(rw http.ResponseWriter, r *http.Request) {
	accessDeniedUrl := "AccessForbidden.html"
	loggedIn := CheckCookies(r, nil)
	if(!loggedIn) {
      http.Redirect(rw, r, accessDeniedUrl, http.StatusSeeOther)
  } else {
      usernameCookie, _ := r.Cookie("username")
      username := usernameCookie.Value
      tmpl := template.Must(template.ParseFiles("../HTML/Links.html"))
      data := LinkData{
          Links: GetLinksForUser(username),
      }
  tmpl.Execute(rw, data)
	}

}
