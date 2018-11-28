package main

import (
	"html/template"
	"net/http"
	"io/ioutil"
	"fmt"
)

var allowedRoles = []string{
	"ADMIN",
}

type UserRole struct {
	RoleName string
	HasAccess  bool
}

type RolePageData struct {
	Username string
	UserRoles []UserRole
}

func ChangeUserRole(rw http.ResponseWriter, r *http.Request) {
	accessDeniedUrl := "AccessForbidden.html"
	loggedIn := CheckCookies(r, allowedRoles)
	if(!loggedIn) {
        http.Redirect(rw, r, accessDeniedUrl, http.StatusSeeOther)
    } else {
		var chosenRoleNames []string
		if err := r.ParseForm(); err != nil {
		   fmt.Println("Error")
		   checkErr(err)
		}
		fmt.Println("Values: ")
		roles := GetAllRoles()
		for _, roleName := range roles {
			if(len(r.Form[roleName]) > 0){
				fmt.Println(roleName)
				chosenRoleNames = append(chosenRoleNames, roleName)
			}
	    }
		fmt.Println(r.Form["username"])
		AdjustUserRoles(r.Form["username"][0], chosenRoleNames)
		http.Redirect(rw, r, "ChooseUser.html", http.StatusSeeOther)
	}

}

func ChooseUser(rw http.ResponseWriter, r *http.Request) {
	accessDeniedUrl := "AccessForbidden.html"
	loggedIn := CheckCookies(r, allowedRoles)
	if(!loggedIn) {
        http.Redirect(rw, r, accessDeniedUrl, http.StatusSeeOther)
    } else {
		dat, err := ioutil.ReadFile("../HTML/ChooseUser.html")
	    checkErr(err)
	    rw.Write(dat)
	}
}

func LoadChangeUserRole(rw http.ResponseWriter, r *http.Request) {

    loggedIn := CheckCookies(r, allowedRoles)
    accessDeniedUrl := "AccessForbidden.html"
    if(!loggedIn) {
        http.Redirect(rw, r, accessDeniedUrl, http.StatusSeeOther)
    }

	r.ParseForm()
	userToChange:= r.FormValue("username")


    tmpl := template.Must(template.ParseFiles("../HTML/ChangeUserRole.html"))
    data := RolePageData{
			Username: userToChange,
			UserRoles: loadUserRoles(userToChange),
		}
	tmpl.Execute(rw, data)
}

func loadUserRoles(username string)([]UserRole){
	return GetUserRolesForUser(username);

}
