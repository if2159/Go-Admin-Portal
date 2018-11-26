
package main

import (
    "net/http"
    "time"
    "io/ioutil"
)

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
    loggedIn := CheckCookies(r)
    newUrl := "links.html"
    if(loggedIn) {
    http.Redirect(rw, r, newUrl, http.StatusSeeOther)
    }
    dat, err := ioutil.ReadFile("../HTML/login.html")
    checkErr(err)
    rw.Write(dat)
}


func CheckCookies(r *http.Request) (bool) {
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

    return CheckSessionId(username, sessionId)
}


func createNewSession(username string, rw http.ResponseWriter){
    sessionId := createSessionId(username)

    expiration := time.Now().Add(1 * time.Hour)
    usernameCookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
    http.SetCookie(rw, &usernameCookie)

    sessionIdCookie := http.Cookie{Name: "sessionId", Value: sessionId, Expires: expiration, HttpOnly:true}

    http.SetCookie(rw, &sessionIdCookie)

}
