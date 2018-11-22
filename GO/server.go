
package main

import (
    "net/http"
    "time"
    "fmt"
    "path/filepath"
)

func main() {
    http.HandleFunc("/Login", Login)
    absPath, _ := filepath.Abs("../HTML")
    http.Handle("/", http.FileServer(http.Dir(absPath)))
    http.ListenAndServe(":80", nil)
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


func createNewSession(username string, rw http.ResponseWriter){
    fmt.Println("1")
    sessionId := createSessionId(username)
    fmt.Println("2")

    expiration := time.Now().Add(365 * 24 * time.Hour)
    usernameCookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
    //fmt.Println(usernameCookie)
    http.SetCookie(rw, &usernameCookie)
    fmt.Println("3")

    sessionIdCookie := http.Cookie{Name: "sessionId", Value: sessionId, Expires: expiration, HttpOnly:true}

    http.SetCookie(rw, &sessionIdCookie)
    //fmt.Println(sessionIdCookie)
    //fmt.Println("4")

}
