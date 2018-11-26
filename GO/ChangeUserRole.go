package main

import (
	"html/template"
	"net/http"
)

type UserRole struct {
	RoleName string
	HasAccess  bool
}

type RolePageData struct {
	PageTitle string
	UserRoles []UserRole
}

func LoadChangeUserRole(rw http.ResponseWriter, r *http.Request) {

    loggedIn := CheckCookies(r)
    accessDeniedUrl := "AccessForbidden.html"
    if(!loggedIn) {
        http.Redirect(rw, r, newUrl, http.StatusSeeOther)
    }

    tmpl := template.Must(template.ParseFiles("../HTML/ChangeUserRole.html"))
    data := RolePageData{
			PageTitle: "Roles",
			UserRoles: loadUserRoles(),
		}
	tmpl.Execute(rw, data)
}

func loadUserRoles(username string)([]UserRole){


}
