
package main

import (
    "html/template"
	"net/http"
    "io/ioutil"
    "fmt"
)


type LinkObj struct {
	LinkName string
}

type RemoveLinkPageData struct {
	PageName string
	Links []LinkObj
}


func LoadRemoveLinks(rw http.ResponseWriter, r *http.Request) {
        loggedIn := CheckCookies(r, allowedRoles)
        accessDeniedUrl := "AccessForbidden.html"
        if(!loggedIn) {
            http.Redirect(rw, r, accessDeniedUrl, http.StatusSeeOther)
        }


        tmpl := template.Must(template.ParseFiles("../HTML/RemoveLinks.html"))

        allLinks := GetAllLinks()
        allLinkObj := stringsToLinkObj(allLinks)

        data := RemoveLinkPageData{
                PageName: "",
                Links: allLinkObj,
            }
        tmpl.Execute(rw, data)
}

func stringsToLinkObj(allLinks []string)([]LinkObj) {
    var objs []LinkObj
    for _, link := range allLinks{
        linkob := LinkObj{
            LinkName:link,
        }
        objs = append(objs, linkob)
    }
    return objs
}

func RemoveLinks(rw http.ResponseWriter, r *http.Request) {
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
		AdjustLinks(chosenLinks)
		http.Redirect(rw, r, "RemoveLinks.html", http.StatusSeeOther)
	}

}

func LoadAddLink(rw http.ResponseWriter, r *http.Request) {
    	accessDeniedUrl := "AccessForbidden.html"
    	loggedIn := CheckCookies(r, allowedRoles)
    	if(!loggedIn) {
            http.Redirect(rw, r, accessDeniedUrl, http.StatusSeeOther)
        } else {
    		dat, err := ioutil.ReadFile("../HTML/AddLink.html")
    	    checkErr(err)
    	    rw.Write(dat)
    	}
    }

func CreateLink(rw http.ResponseWriter, r *http.Request) {
    linkName := r.FormValue("linkName")
    linkUrl := r.FormValue("linkUrl")

    if(linkName != "" && linkUrl != ""){
        InsertNewLink(linkName, linkUrl)
    }

    chooseRoleURL := "AddLink.html"
    http.Redirect(rw, r, chooseRoleURL, http.StatusSeeOther)

}
