
package main

import (
    "net/http"
    "time"
    "io/ioutil"
)
 var availableRoles = []string{
                            "GUEST",
                            "ADMIN",
                            "USER",
                        }

func Login(rw http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    password := r.FormValue("password")

    var resp []byte

    if(comparePassword(username, []byte(password))) {
        createNewSession(username, rw);
        resp = []byte("Success!")
    } else{
        resp = []byte("Fail")
    }

    rw.Write(resp)
}


func LoadLogin(rw http.ResponseWriter, r *http.Request) {
    loggedIn := CheckCookies(r, availableRoles)
    newUrl := "Links.html"
    if(loggedIn) {
    http.Redirect(rw, r, newUrl, http.StatusSeeOther)
    }
    dat, err := ioutil.ReadFile("../HTML/login.html")
    checkErr(err)
    rw.Write(dat)
}



func CheckCookies(r *http.Request, allowedRoles []string) (bool) {
    usernameCookie, _ := r.Cookie("username")
    if(usernameCookie == nil){
        return false
    }
    username := usernameCookie.Value

    sessionIdCookie, _ := r.Cookie("sessionId")
    if(sessionIdCookie == nil){
        return false
    }
    sessionId := sessionIdCookie.Value

    if(CheckSessionId(username, sessionId)){
        if(allowedRoles != nil){
            var roles = GetRolesForUser(username)
            for _, role := range roles{
                if(stringInSlice(role, allowedRoles)){
                    return true
                }
            }
        }else{
            return true;
        }
    }
    return false
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}


func createNewSession(username string, rw http.ResponseWriter){
    sessionId := createSessionId(username)

    expiration := time.Now().Add(1 * time.Hour)
    usernameCookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
    http.SetCookie(rw, &usernameCookie)

    sessionIdCookie := http.Cookie{Name: "sessionId", Value: sessionId, Expires: expiration, HttpOnly:true}

    http.SetCookie(rw, &sessionIdCookie)

}
