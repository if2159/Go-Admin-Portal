package main

import (
    _ "github.com/go-sql-driver/mysql"
    "golang.org/x/crypto/bcrypt"
    "database/sql"
    "strings"
    "fmt"
)

var con *sql.DB

type User struct {
    uid string
    username string
    roleId int
    hashedPass []byte
}


func getStoredUser(username string) (User){
    startConnection()
    rows, err := con.Query("SELECT UID, USERNAME, ROLE_ID, HASHED_PASS FROM USERS WHERE USERNAME = ?", username)
    var user User
    for rows.Next() {
        var uid string
        var username string
        var roleId int
        var hashedPass []byte
        err = rows.Scan(&uid, &username, &roleId, &hashedPass)
        checkErr(err)
        user = User{uid, username, roleId, hashedPass}
    }

    return user
}

func startConnection(){
    fmt.Println("start")
    if(con == nil){
        user, password, host := loadDbProperties()

        conn, err := sql.Open("mysql", user + ":@" + password + "@tcp(" + host + ":3306)/portal?charset=utf8")
        con = conn
        checkErr(err)
    }
    fmt.Println("startEnd")
}

func CheckSessionId(username, sessionId string) (bool){
    startConnection()
    rows, err := con.Query("SELECT SESSION_ID, UID FROM portal.SESSIONS WHERE UID=(SELECT UID FROM USERS WHERE USERNAME=?) AND SESSION_ID=?", username, sessionId);
    checkErr(err)

    for rows.Next() {
        var sessionId string
        var uid string
        err = rows.Scan(&sessionId, &uid)
        checkErr(err)

        return true
    }
    return false
}

func updatePassword(){

    pass := []byte("password123")

   // Hashing the password with the default cost of 10
    hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)

    con.Query("UPDATE USERS SET HASHED_PASS = ? WHERE USERNAME=?", hashedPassword, "CHRISTIAN")
    checkErr(err)
}


func comparePassword(username string, password []byte) (bool){
    user := getStoredUser(strings.ToUpper(username))
    updatePassword()
    storedHash := user.hashedPass
    return (bcrypt.CompareHashAndPassword(storedHash, password) == nil)
}

func createSessionId(username string) (string){
    fmt.Println("Creating Session ID")
    startConnection()

    con.Query("INSERT INTO SESSIONS(UID, SESSION_ID) VALUES((SELECT UID FROM USERS WHERE USERNAME=?), UUID())", username)
    rows, err := con.Query("SELECT SESSION_ID, UID FROM portal.SESSIONS WHERE UID=(SELECT UID FROM USERS WHERE USERNAME=?) AND LOGIN_TIME = (SELECT MAX(LOGIN_TIME) FROM SESSIONS WHERE UID = (SELECT UID FROM USERS WHERE USERNAME=?))", username, username);
    checkErr(err)

    for rows.Next() {
        var sessionId string
        var uid string
        err = rows.Scan(&sessionId, &uid)
        checkErr(err)

        return sessionId
    }
    return ""
}


func checkErr(err error) {

    if err != nil {
        panic(err)
    }
}
