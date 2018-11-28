package main

import (
	"html/template"
	"net/http"
	"io/ioutil"
	"fmt"
)

type RoleLink struct {
	LinkName string
	HasAccess  bool
}

type LinkPageData struct {
	RoleName string
	RoleLinks []RoleLink
}

func ChangeRoleLinks(rw http.ResponseWriter, r *http.Request) {
	accessDeniedUrl := "AccessForbidden.html"
	loggedIn := CheckCookies(r, allowedRoles)
	if(!loggedIn) {
        http.Redirect(rw, r, accessDeniedUrl, http.StatusSeeOther)
    } else {
		var chosenLinks []string
		if err := r.ParseForm(); err != nil {
		   fmt.Println("Error")
		   checkErr(err)
		}
		fmt.Println("Values: ")
		links := GetAllLinks()
		for _, linkName := range links {
			if(len(r.Form[linkName]) > 0){
				fmt.Println(linkName)
				chosenLinks = append(chosenLinks, linkName)
			}
	    }
		fmt.Println(r.Form["roleName"])
		AdjustRoleLinks(r.Form["roleName"][0], chosenLinks)
		http.Redirect(rw, r, "ChooseRole.html", http.StatusSeeOther)
	}

}

func ChooseRole(rw http.ResponseWriter, r *http.Request) {
	accessDeniedUrl := "AccessForbidden.html"
	loggedIn := CheckCookies(r, allowedRoles)
	if(!loggedIn) {
        http.Redirect(rw, r, accessDeniedUrl, http.StatusSeeOther)
    } else {
		dat, err := ioutil.ReadFile("../HTML/ChooseRole.html")
	    checkErr(err)
	    rw.Write(dat)
	}
}

func LoadChangeRoleLinks(rw http.ResponseWriter, r *http.Request) {

    loggedIn := CheckCookies(r, allowedRoles)
    accessDeniedUrl := "AccessForbidden.html"
    if(!loggedIn) {
        http.Redirect(rw, r, accessDeniedUrl, http.StatusSeeOther)
    }

	r.ParseForm()
	userToChange:= r.FormValue("roleName")


    tmpl := template.Must(template.ParseFiles("../HTML/ChangeRoleLinks.html"))
    data := LinkPageData{
			RoleName: userToChange,
			RoleLinks: loadRoleLinks(userToChange),
		}
	tmpl.Execute(rw, data)
}

func loadRoleLinks(roleName string)([]RoleLink){
	return GetRoleLinkForRole(roleName);

}
