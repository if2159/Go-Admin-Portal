
package main

import (
    "net/http"
    "path/filepath"
    "fmt"
)

func main() {
    fmt.Println("Server Starting...")
    port := ":80"
    fmt.Println("Server started on port: " + port)
    http.HandleFunc("/Login", Login)
    absPath, _ := filepath.Abs("../HTML")
    http.HandleFunc("/login.html", LoadLogin)
    http.HandleFunc("/ChooseUser.html", ChooseUser)
    http.HandleFunc("/ChangeUserRole", ChangeUserRole)
    http.HandleFunc("/ChangeUserRole.html", LoadChangeUserRole)
    http.Handle("/", http.FileServer(http.Dir(absPath)))
    http.ListenAndServe(port, nil)
}
